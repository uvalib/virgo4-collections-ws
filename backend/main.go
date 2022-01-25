package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "2.3.1"

func main() {
	log.Printf("===> Collections Context service starting up <===")

	// Get config params and use them to init service context. Any issues are fatal
	cfg := LoadConfiguration()
	svc := InitializeService(version, cfg)

	log.Printf("INFO: setup routes...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	router.Use(cors.New(config))

	router.GET("/favicon.ico", svc.ignoreFavicon)
	router.GET("/version", svc.getVersion)
	router.GET("/healthcheck", svc.healthCheck)
	api := router.Group("/api")
	{
		api.GET("/logos", svc.userMiddleware, svc.getLogos)
		api.GET("/lookup", svc.lookupCollectionContext)
		api.GET("/features", svc.userMiddleware, svc.getFeatures)
		api.GET("/collections", svc.getCollections)
		api.POST("/collections", svc.userMiddleware, svc.addOrUpdateCollection)
		api.POST("/collections/:id/logo", svc.userMiddleware, svc.uploadLogo)
		api.DELETE("/collections/:id/logo/:fn", svc.userMiddleware, svc.deletePendingLogo)
		api.DELETE("/collections/:id", svc.userMiddleware, svc.deleteCollection)
		api.GET("/collections/:id", svc.userMiddleware, svc.getCollectionDetails)
		api.GET("/collections/:id/dates", svc.collectionMiddleware, svc.getCollectioDates)
		api.GET("/collections/:id/items/:date/next", svc.collectionMiddleware, svc.getNextItem)
		api.GET("/collections/:id/items/:date/previous", svc.collectionMiddleware, svc.getPreviousItem)
	}

	// Note: in dev mode, this is never actually used. The front end is served
	// by yarn and it proxies all requests to the API to the routes above
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("INFO: start service v%s on port %s", version, portStr)
	log.Fatal(router.Run(portStr))
}
