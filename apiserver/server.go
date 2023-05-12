package apiserver

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/aescanero/openldap-node/config"
	"github.com/gin-gonic/gin"
)

//go:embed dashboard
var dashboard embed.FS

func Server(apiconfig config.Config) {
	serverRoot, err := fs.Sub(dashboard, "dashboard")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.New()
	router.StaticFS("/dashboard", http.FileSystem(http.FS(serverRoot)))
	router.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/auth"
		router.HandleContext(c)
	})
	router.GET("/api/hello", hello)
	router.GET("/api/monitor", func(c *gin.Context) {
		monitor(c, apiconfig)
	})
	router.Run(":9090")
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "world"})
}

func auth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "world"})
}
