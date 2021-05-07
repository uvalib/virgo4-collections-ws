package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ServiceContext contains common data used by all handlers
type ServiceContext struct {
	Version    string
	Solr       SolrConfig
	HTTPClient *http.Client
}

// InitializeService sets up the service context for all API handlers
func InitializeService(version string, cfg *ServiceConfig) *ServiceContext {
	ctx := ServiceContext{Version: version, Solr: cfg.Solr}

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
	}
	hcMap := make(map[string]hcResp)
	hcMap["collectionsvc"] = hcResp{Healthy: true}

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
