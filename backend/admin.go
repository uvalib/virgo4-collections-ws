package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
