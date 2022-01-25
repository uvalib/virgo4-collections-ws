package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type solrResponseHeader struct {
	Status int `json:"status,omitempty"`
}

type solrDocument map[string]interface{}

type solrResponseDocuments struct {
	NumFound int            `json:"numFound,omitempty"`
	Start    int            `json:"start,omitempty"`
	Docs     []solrDocument `json:"docs,omitempty"`
}

type solrResponse struct {
	Header   solrResponseHeader    `json:"responseHeader,omitempty"`
	Response solrResponseDocuments `json:"response,omitempty"`
}

type feature struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type collection struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ItemLabel   string    `json:"item_label,omitempty"`
	FilterName  string    `json:"filter_name"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitempty"`
	Active      bool      `json:"active,omitempty"`
	Features    []feature `json:"features,omitempty" gorm:"many2many:collection_features"`
	Image       *image    `json:"image,omitempty"`
}

type image struct {
	ID           int64  `json:"id"`
	CollectionID int64  `json:"-"`
	AltText      string `json:"alt_text,omitempty"`
	Title        string `json:"title,omitempty"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Filename     string `json:"filename"`
	URL          string `json:"url"`
}

func (svc *ServiceContext) getCollections(c *gin.Context) {
	all := c.Query("all")
	log.Printf("INFO: get a list of collections")

	var recs []collection
	var dbResp *gorm.DB
	if all != "" {
		log.Printf("INFO: get all collections, including inactive")
		dbResp = svc.GDB.Select("id,title,filter_name").Order("title asc").Find(&recs)
	} else {
		log.Printf("INFO: get all active collections")
		dbResp = svc.GDB.Select("id,title,filter_name").Order("title asc").Where("active=?", true).Find(&recs)
	}

	if dbResp.Error != nil {
		log.Printf("ERROR: unable to get collections: %s", dbResp.Error.Error())
		c.String(http.StatusInternalServerError, dbResp.Error.Error())
		return
	}

	c.JSON(http.StatusOK, recs)
}

func (svc *ServiceContext) lookupCollectionContext(c *gin.Context) {
	rawName := c.Query("q")
	log.Printf("INFO: lookup collection context for [%s]", rawName)

	var rec collection
	dbResp := svc.GDB.Preload("Image").Preload("Features").Where("title=? and active=?", rawName, true).First(&rec)
	if dbResp.Error != nil {
		if errors.Is(dbResp.Error, gorm.ErrRecordNotFound) {
			log.Printf("INFO: no collection context found for %s", rawName)
			c.String(http.StatusNotFound, "not found")
		} else {
			log.Printf("ERROR: context lookup for %s failed: %s", rawName, dbResp.Error.Error())
			c.String(http.StatusInternalServerError, dbResp.Error.Error())
		}
		return
	}

	if rec.Image != nil {
		rec.Image.URL = fmt.Sprintf("%s/%s", svc.BaseImageURL, rec.Image.Filename)
	}

	c.JSON(http.StatusOK, rec)
}

// collection middleware accepts a collection ID and finds the associated filter value (title)
func (svc *ServiceContext) collectionMiddleware(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("INFO: lookup collection id %d", id)

	var coll collection
	dbResp := svc.GDB.Select("title").Where("id=? and active=?", id, true).First(&coll)
	if dbResp.Error != nil {
		if errors.Is(dbResp.Error, gorm.ErrRecordNotFound) {
			log.Printf("INFO: collection %d not found", id)
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Printf("ERROR: get collection %d failed: %s", id, dbResp.Error.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Set("key", coll.Title)
	log.Printf("INFO collection id %d=[%s]", id, coll.Title)
	c.Next()
}

func (svc *ServiceContext) getCollectioDates(c *gin.Context) {
	key := c.GetString("key")
	year := c.Query("year")
	if year == "" {
		log.Printf("ERROR: year is required")
		c.String(http.StatusBadRequest, "year param is required")
		return
	}
	if len(year) != 4 {
		log.Printf("ERROR: year %s is invalid", year)
		c.String(http.StatusBadRequest, "year param must be of the format YYYY")
		return
	}
	_, err := strconv.Atoi(year)
	if err != nil {
		log.Printf("ERROR: year %s is invalid: %s", year, err.Error())
		c.String(http.StatusBadRequest, "year param must be of the format YYYY")
		return
	}
	log.Printf("INFO: get all item publication dates for %s in %s", year, key)

	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id,published_date")
	qParams = append(qParams, "rows=365")
	qParams = append(qParams, "start=0")
	qParams = append(qParams, "sort=published_date%20asc")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[%s-01-01T00:00:00.000Z TO %s-12-31T00:00:00.000Z]",
		key, year, year)
	qParams = append(qParams, fmt.Sprintf("q=%s", url.QueryEscape(q)))
	solrURL := fmt.Sprintf("%s?%s", svc.Solr.selectURL(), strings.Join(qParams, "&"))
	resp, err := svc.getAPIResponse(solrURL)
	if err != nil {
		log.Printf("ERROR: solr query failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var respJSON solrResponse
	err = json.Unmarshal(resp, &respJSON)
	if err != nil {
		log.Printf("ERROR: unable to parse solr response: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type pubDate struct {
		Date string `json:"date"`
		PID  string `json:"pid"`
	}
	log.Printf("%+v", respJSON.Response.Docs)
	out := make([]pubDate, 0)
	for _, doc := range respJSON.Response.Docs {
		if dateStr, ok := doc["published_date"].(string); ok {
			cleanDate := strings.Split(dateStr, "T")[0]
			item := pubDate{Date: cleanDate, PID: doc["id"].(string)}
			out = append(out, item)
		}
	}

	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) getPreviousItem(c *gin.Context) {
	key := c.GetString("key")
	date := c.Param("date")
	log.Printf("INFO: navigate to item after %s in %s", date, key)

	// subtract 1 milisecond from the passed time; convert to time, then subtract,
	// then back to string
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)
	old := t.Add(-time.Millisecond)
	date = old.Format(time.RFC3339)

	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[* TO %s]", key, date)
	qParams = append(qParams, fmt.Sprintf("q=%s", url.QueryEscape(q)))
	qParams = append(qParams, "rows=1")
	qParams = append(qParams, "start=0")
	qParams = append(qParams, "sort=published_date%20desc")
	solrURL := fmt.Sprintf("%s?%s", svc.Solr.selectURL(), strings.Join(qParams, "&"))
	resp, err := svc.getAPIResponse(solrURL)
	if err != nil {
		log.Printf("ERROR: solr query failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var respJSON solrResponse
	err = json.Unmarshal(resp, &respJSON)
	if err != nil {
		log.Printf("ERROR: unable to parse solr response: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	firstHit := respJSON.Response.Docs[0]
	PID := firstHit["id"].(string)
	log.Printf("INFO: item before %s is %s", date, PID)
	c.String(http.StatusOK, PID)
}

func (svc *ServiceContext) getNextItem(c *gin.Context) {
	key := c.GetString("key")
	date := c.Param("date")
	log.Printf("INFO: navigate to item before %s in %s", date, key)
	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[%sT00:00:00.001Z TO *]", key, date)
	qParams = append(qParams, fmt.Sprintf("q=%s", url.QueryEscape(q)))
	qParams = append(qParams, "rows=1")
	qParams = append(qParams, "start=0")
	qParams = append(qParams, "sort=published_date%20asc")
	solrURL := fmt.Sprintf("%s?%s", svc.Solr.selectURL(), strings.Join(qParams, "&"))
	resp, err := svc.getAPIResponse(solrURL)
	if err != nil {
		log.Printf("ERROR: solr query failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var respJSON solrResponse
	err = json.Unmarshal(resp, &respJSON)
	if err != nil {
		log.Printf("ERROR: unable to parse solr response: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	firstHit := respJSON.Response.Docs[0]
	PID := firstHit["id"].(string)
	log.Printf("INFO: item after %s is %s", date, PID)
	c.String(http.StatusOK, PID)
}
