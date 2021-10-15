package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (svc *ServiceContext) userMiddleware(c *gin.Context) {
	log.Printf("INFO: checking authentication headers")
	log.Printf("INFO: start header dump ========================================")
	for name, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("%s=%s\n", name, value)
		}
	}
	log.Printf("INFO: end header dump ==========================================")

	computingID := c.GetHeader("remote_user")
	if computingID == "" {
		log.Printf("INFO: missing netbadge user header")
		computingID = "ANONYMOUS"
	}

	c.Set("user", computingID)
	c.Next()
}
