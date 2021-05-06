package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (svc *ServiceContext) getCollectionContext(c *gin.Context) {
	rawName := c.Param("name")
	log.Printf("Get collection context for [%s]", rawName)

	if !svc.DevMode {
		v4HostHeader := c.Request.Header.Get("V4Host")
		log.Printf("Request V4Host header: %s", v4HostHeader)
		if strings.Index(v4HostHeader, "-dev") > -1 {
			log.Printf("This request is from a dev server")
		} else {
			log.Printf("INFO: this request is from a non-dev server; return blank response")
			var empty interface{}
			c.JSON(http.StatusOK, empty)
			return
		}
	}

	name := strings.ToLower(rawName)
	name = strings.ReplaceAll(name, " ", "_")
	fileName := fmt.Sprintf("./data/%s.json", name)
	if _, err := os.Stat(fileName); err == nil {
		log.Printf("INFO: collection context exists for %s (%s)", rawName, name)
		c.Header("Content-Type", "application/json")
		c.File(fileName)
	} else {
		log.Printf("WARNING: no collection context found for %s", name)
		c.String(http.StatusNotFound, "not found")
	}
}
