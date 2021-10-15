package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

// ServiceContext contains common data used by all handlers
type ServiceContext struct {
	Version       string
	BaseImageURL  string
	Solr          SolrConfig
	HTTPClient    *http.Client
	DB            *dbx.DB
	S3ImageBucket string
	S3Uploader    *s3manager.Uploader
	S3Service     *s3.S3
}

// InitializeService sets up the service context for all API handlers
func InitializeService(version string, cfg *ServiceConfig) *ServiceContext {
	ctx := ServiceContext{Version: version, Solr: cfg.Solr, BaseImageURL: cfg.ImageBaseURL}

	log.Printf("INFO: init S3 session and uploader")
	bName := strings.Split(cfg.ImageBaseURL, "https://")[1]
	bName = strings.Split(bName, ".")[0]
	ctx.S3ImageBucket = bName
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	ctx.S3Uploader = s3manager.NewUploader(sess)
	ctx.S3Service = s3.New(sess)
	log.Printf("INFO: S3 bucket %s and upload manager initailized", ctx.S3ImageBucket)

	log.Printf("Connect to Postgres")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Host, cfg.DB.Port)
	db, err := dbx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.LogFunc = log.Printf
	ctx.DB = db

	log.Printf("INFO: create HTTP client...")
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 600 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	ctx.HTTPClient = &http.Client{
		Transport: defaultTransport,
		Timeout:   5 * time.Second,
	}
	log.Printf("INFO: HTTP Client created")

	return &ctx
}

func (svc *ServiceContext) ignoreFavicon(c *gin.Context) {
}

func (svc *ServiceContext) getVersion(c *gin.Context) {
	build := "unknown"
	// working directory is the bin directory, and build tag is in the root
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = svc.Version
	vMap["build"] = build
	c.JSON(http.StatusOK, vMap)
}

func (svc *ServiceContext) healthCheck(c *gin.Context) {
	type hcResp struct {
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
		Version int    `json:"version,omitempty"`
	}
	hcMap := make(map[string]hcResp)
	hcMap["collectionsvc"] = hcResp{Healthy: true}

	tq := svc.DB.NewQuery("select * from schema_migrations order by version desc limit 1")
	var schema struct {
		Version int  `db:"version"`
		Dirty   bool `db:"dirty"`
	}
	err := tq.One(&schema)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		hcMap["postgres"] = hcResp{Healthy: false, Message: err.Error()}
	} else {
		log.Printf("Schema info - Version: %d, Dirty: %t", schema.Version, schema.Dirty)
		if schema.Dirty {
			hcMap["postgres"] = hcResp{Healthy: false, Message: fmt.Sprintf("Schema %d is marked dirty", schema.Version)}
		} else {
			// check that the highest numbered migration matches DB version value
			latest := getLatestMigrationNumber()
			if latest > 0 && latest != schema.Version {
				hcMap["postgres"] = hcResp{Healthy: false, Message: fmt.Sprintf("Schema out of date. Reported version: %d, latest: %d", schema.Version, latest)}
			}
		}
		hcMap["postgres"] = hcResp{Healthy: true, Version: schema.Version}
	}

	c.JSON(http.StatusOK, hcMap)
}

func (svc *ServiceContext) getAPIResponse(url string) ([]byte, error) {
	log.Printf("INFO: GET API Response from %s, timeout  %.0f sec", url, svc.HTTPClient.Timeout.Seconds())
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36")

	startTime := time.Now()
	resp, rawErr := svc.HTTPClient.Do(req)
	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)
	bodyBytes, err := handleAPIResponse(url, resp, rawErr)
	if err != nil {
		log.Printf("ERROR: %s : %s. Elapsed Time: %d (ms)", url, err.Error(), elapsedMS)
		return nil, err
	}

	log.Printf("INFO: successful response from %s. Elapsed Time: %d (ms)", url, elapsedMS)
	return bodyBytes, nil
}

func handleAPIResponse(url string, resp *http.Response, rawErr error) ([]byte, error) {
	if rawErr != nil {
		status := http.StatusBadRequest
		errMsg := rawErr.Error()
		if strings.Contains(rawErr.Error(), "Timeout") {
			status = http.StatusRequestTimeout
			errMsg = fmt.Sprintf("%s timed out", url)
		} else if strings.Contains(rawErr.Error(), "connection refused") {
			status = http.StatusServiceUnavailable
			errMsg = fmt.Sprintf("%s refused connection", url)
		}
		err := fmt.Errorf("%d: %s", status, errMsg)
		return nil, err
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		status := resp.StatusCode
		errMsg := string(bodyBytes)
		err := fmt.Errorf("%d: %s", status, errMsg)
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, nil
}

func getLatestMigrationNumber() int {
	// on deployed systems, the migrations can be found in ../db/*.sql (we run from ./bin)
	// on dev systems run with 'go run', they are in ./backend/db/migrations/*.sql
	// try both; if files found return the latest number. If none found, just return 0.
	tgts := []string{"../db", "./db/migrations"}
	migrateDir := ""
	for _, dir := range tgts {
		_, err := os.Stat(dir)
		if err == nil {
			migrateDir = dir
			break
		}
	}

	if migrateDir == "" {
		log.Printf("WARN: DB Migration directory not found")
		return 0
	}

	files, err := ioutil.ReadDir(migrateDir)
	if err != nil {
		log.Printf("WARN: DB Migration directory unreadable: %s", err.Error())
		return 0
	}

	maxNum := -1
	maxFile := ""
	for _, f := range files {
		fname := f.Name()
		if strings.Contains(fname, "up.sql") {
			numStr := strings.Split(fname, "_")[0]
			num, _ := strconv.Atoi(numStr)
			if num > maxNum {
				maxNum = num
				maxFile = fname
			}
		}
	}

	// there are up/down files for each migration
	log.Printf("Last migration file found: %s, version: %d", maxFile, maxNum)
	return maxNum
}
