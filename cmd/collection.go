package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func (svc *ServiceContext) getCollectionContext(c *gin.Context) {
	rawName := c.Param("name")
	log.Printf("Get collection context for [%s]", rawName)

	fileName := collectionFilename(rawName)
	if _, err := os.Stat(fileName); err == nil {
		c.Header("Content-Type", "application/json")
		c.File(fileName)
	} else {
		log.Printf("WARNING: no collection context found for %s", fileName)
		c.String(http.StatusNotFound, "not found")
	}
}

func collectionFilename(rawName string) string {
	name := strings.ToLower(rawName)
	name = strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("./data/%s.json", name)
}

func (svc *ServiceContext) getCollectioDates(c *gin.Context) {
	name := c.Param("name")
	year := c.Query("year")
	if year == "" {
		log.Printf("ERROR: year is required")
		c.String(http.StatusBadRequest, "year param is required")
		return
	}
	log.Printf("INFO: get all item publication dates for %s in %s", year, name)

	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id,published_date")
	qParams = append(qParams, "rows=365")
	qParams = append(qParams, "start=0")
	qParams = append(qParams, "sort=published_date%20asc")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[%s-01-01T00:00:00.000Z TO %s-12-31T00:00:00.000Z]",
		name, year, year)
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
	name := c.Param("name")
	date := c.Param("date")
	log.Printf("INFO: navigate to item after %s in %s", date, name)

	// subtract 1 milisecond from the passed time; convert to time, then subtract,
	// then back to string
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)
	old := t.Add(-time.Millisecond)
	date = old.Format(time.RFC3339)

	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[* TO %s]", name, date)
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
	name := c.Param("name")
	date := c.Param("date")
	log.Printf("INFO: navigate to item before %s in %s", date, name)
	qParams := make([]string, 0)
	qParams = append(qParams, "fl=id")
	q := fmt.Sprintf("digital_collection_f:\"%s\"+published_date:[%sT00:00:00.001Z TO *]", name, date)
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
