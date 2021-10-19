package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
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

type collectionRec struct {
	ID          int64          `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	ItemLabel   string         `db:"item_label"`
	FilterName  string         `db:"filter_name"`
	StartDate   sql.NullString `db:"start_date"`
	EndDate     sql.NullString `db:"end_date"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (c collectionRec) TableName() string {
	return "collections"
}

type imageRec struct {
	ID           int64          `db:"id"`
	CollectionID int64          `db:"collection_id"`
	AltText      sql.NullString `db:"alt_text"`
	Title        sql.NullString `db:"title"`
	Width        int            `db:"width"`
	Height       int            `db:"height"`
	Filename     string         `db:"filename"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (c imageRec) TableName() string {
	return "images"
}

type imageJSON struct {
	AltText string `json:"alt_text,omitempty"`
	Title   string `json:"title,omitempty"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	URL     string `json:"url"`
}

type collectionJSON struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ItemLabel   string     `json:"item_label"`
	FilterName  string     `json:"filter_name"`
	StartDate   string     `json:"start_date,omitempty"`
	EndDate     string     `json:"end_date,omitempty"`
	Features    []string   `json:"features"`
	Image       *imageJSON `json:"image,omitempty"`
}

func (coll *collectionJSON) getFeatures(db *dbx.DB) {
	log.Printf("INFO: get collection [%s] features", coll.Title)
	q := db.NewQuery("select f.name from features f inner join collection_features cf on cf.feature_id = f.id where cf.collection_id={:cid}")
	q.Bind(dbx.Params{"cid": coll.ID})
	rows, err := q.Rows()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: no features found for collection [%s]", coll.Title)
		} else {
			log.Printf("ERROR: unable to lookup features for [%s]: %s", coll.Title, err.Error())
		}
	} else {
		for rows.Next() {
			val := ""
			rows.Scan(&val)
			coll.Features = append(coll.Features, val)
		}
	}
}

func (coll *collectionJSON) getImages(db *dbx.DB, baseImageURL string) {
	log.Printf("INFO: get collection [%s] images", coll.Title)
	var image imageRec
	q := db.NewQuery("select * from images where collection_id={:cid} limit 1")
	q.Bind(dbx.Params{"cid": coll.ID})
	err := q.One(&image)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: no images found for collection [%s]", coll.Title)
		} else {
			log.Printf("ERROR: unable to lookup images for [%s]: %s", coll.Title, err.Error())
		}
	} else {
		imgJSON := imageJSON{Width: image.Width, Height: image.Height}
		if image.AltText.Valid {
			imgJSON.AltText = image.AltText.String
		}
		if image.Title.Valid {
			imgJSON.Title = image.Title.String
		}
		imgJSON.URL = fmt.Sprintf("%s/%s", baseImageURL, image.Filename)
		coll.Image = &imgJSON
	}
}

// initialize JSON collection data from a DB rec
func collectionfromDB(rec collectionRec) *collectionJSON {
	c := collectionJSON{ID: rec.ID, Title: rec.Title, FilterName: rec.FilterName, ItemLabel: rec.ItemLabel, Features: make([]string, 0)}
	if rec.Description.Valid {
		c.Description = rec.Description.String
	}
	if rec.StartDate.Valid {
		c.StartDate = rec.StartDate.String
	}
	if rec.EndDate.Valid {
		c.EndDate = rec.EndDate.String
	}
	return &c
}

func (svc *ServiceContext) getCollections(c *gin.Context) {
	log.Printf("INFO: get a list of collections")
	type collResp struct {
		ID     int64  `db:"id" json:"id"`
		Filter string `db:"filter_name" json:"facet"`
		Title  string `db:"title" json:"title"`
	}
	q := svc.DB.NewQuery("select id,title,filter_name from collections order by title asc")
	var recs []collResp
	err := q.All(&recs)
	if err != nil {
		log.Printf("ERROR: unable to get collections: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recs)
}

func (svc *ServiceContext) lookupCollectionContext(c *gin.Context) {
	rawName := c.Query("q")
	log.Printf("INFO: lookup collection context for [%s]", rawName)

	var rec collectionRec
	q := svc.DB.NewQuery("select * from collections where title={:fv}")
	q.Bind(dbx.Params{"fv": rawName})
	err := q.One(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: no collection context found for %s", rawName)
			c.String(http.StatusNotFound, "not found")
		} else {
			log.Printf("ERROR: contexed lookup for %s failed: %s", rawName, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := collectionfromDB(rec)
	out.getFeatures(svc.DB)
	out.getImages(svc.DB, svc.BaseImageURL)

	c.JSON(http.StatusOK, out)
}

// collection middleware accepts a collection ID and finds the associated filter value
func (svc *ServiceContext) collectionMiddleware(c *gin.Context) {
	id := c.Param("id")
	log.Printf("INFO: lookup collection id %s", id)

	q := svc.DB.NewQuery("select title from collections where id={:id}")
	q.Bind(dbx.Params{"id": id})
	name := ""
	err := q.Row(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: collection %s not found", id)
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Printf("ERROR: get collection %s failed: %s", id, err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Set("key", name)
	log.Printf("INFO collection id %s=[%s]", id, name)
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
