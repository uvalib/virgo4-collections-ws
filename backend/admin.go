package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

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
