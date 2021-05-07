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
