package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type updateCollectionRequest struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ItemLabel     string `json:"itemLabel"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Filter        string `json:"filter"`
	FeatureIDs    []int  `json:"features"`
	ImageFileName string `json:"imageFile"`
	ImageTitle    string `json:"imageTitle"`
	ImageAlt      string `json:"imageAlt"`
}

func (svc *ServiceContext) getFeatures(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of features", user)
	type featureResp struct {
		ID   int64  `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
	var features []featureResp
	q := svc.DB.NewQuery("select id,name from features order by name asc")
	err := q.All(&features)
	if err != nil {
		log.Printf("ERROR: unable to get features: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, features)
}

func (svc *ServiceContext) getCollections(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of collections", user)
	type collResp struct {
		ID    int64  `db:"id" json:"id"`
		Title string `db:"title" json:"title"`
	}
	q := svc.DB.NewQuery("select id,title from collections order by title asc")
	var recs []collResp
	err := q.All(&recs)
	if err != nil {
		log.Printf("ERROR: unable to get collections: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recs)
}

func (svc *ServiceContext) getCollectionDetails(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s is is requesting collection %s details", user, id)

	var rec collectionRec
	q := svc.DB.NewQuery("select * from collections where id={:cid}")
	q.Bind(dbx.Params{"cid": id})
	err := q.One(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: no collection context found for %s", id)
			c.String(http.StatusNotFound, "not found")
		} else {
			log.Printf("ERROR: contexed lookup for %s failed: %s", id, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := collectionfromDB(rec)
	out.getFeatures(svc.DB)
	out.getImages(svc.DB, svc.BaseImageURL)

	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) deleteCollection(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s requests delete of collection %s", user, id)
	q := svc.DB.NewQuery("select * from collections where id={:cid}")
	q.Bind(dbx.Params{"cid": id})
	var rec collectionRec
	err := q.One(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("WARNING: collection %s not found", id)
			c.String(http.StatusNotFound, fmt.Sprintf("%s not found", id))
		} else {
			log.Printf("ERROR: unable to find collection %s: %s", id, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = svc.DB.Model(&rec).Delete()
	if err != nil {
		log.Printf("ERROR: unable to delete ollection %s: %s", id, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "deleted")
}

func (svc *ServiceContext) addOrUpdateCollection(c *gin.Context) {
	user := c.GetString("user")
	var req updateCollectionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("ERROR: %s update request with invalid payload: %v", user, err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	updateRec := collectionRec{ID: int64(req.ID), Title: req.Title, ItemLabel: req.ItemLabel, FilterName: req.Filter}
	if req.Description != "" {
		updateRec.Description.String = req.Description
		updateRec.Description.Valid = true
	}
	if req.StartDate != "" {
		updateRec.StartDate.String = req.StartDate
		updateRec.StartDate.Valid = true
	}
	if req.EndDate != "" {
		updateRec.EndDate.String = req.EndDate
		updateRec.EndDate.Valid = true
	}

	if req.ID == 0 {
		log.Printf("INFO: %s add collection %+v", user, req)
		addErr := svc.DB.Model(&updateRec).Insert()
		if addErr != nil {
			log.Printf("ERROR: %s add %v failed: %s", user, req, addErr.Error())
			c.String(http.StatusInternalServerError, addErr.Error())
			return
		}
	} else {
		log.Printf("INFO: %s update collection %+v", user, req)
		upErr := svc.DB.Model(&updateRec).Exclude("ID").Update()
		if upErr != nil {
			log.Printf("ERROR: %s update %v failed: %s", user, req, upErr.Error())
			c.String(http.StatusInternalServerError, upErr.Error())
			return
		}
		q := svc.DB.NewQuery("delete from collection_features where collection_id={:cid}")
		q.Bind(dbx.Params{"cid": updateRec.ID})
		_, upErr = q.Execute()
		if upErr != nil {
			log.Printf("ERROR: %s reset features failed: %s", user, upErr.Error())
			c.String(http.StatusInternalServerError, upErr.Error())
			return
		}
	}

	log.Printf("INFO: adding features to collection %d", updateRec.ID)
	qStr := "insert into collection_features (collection_id, feature_id) values "
	vals := make([]string, 0)
	for _, featureID := range req.FeatureIDs {
		vals = append(vals, fmt.Sprintf("(%d,%d)", updateRec.ID, featureID))
	}
	qStr += strings.Join(vals, ",")
	q := svc.DB.NewQuery(qStr)
	_, err = q.Execute()
	if err != nil {
		log.Printf("ERROR: %s add features failed: %s", user, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// convert DB structs into JSON response
	out := collectionfromDB(updateRec)
	out.getFeatures(svc.DB)
	out.getImages(svc.DB, svc.BaseImageURL)
	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) uploadLogo(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s is uploading a new logo for collection %s", user, id)

	c.String(http.StatusInternalServerError, "NO")
}
