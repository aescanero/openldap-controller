package apiserver

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aescanero/openldap-node/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed dashboard/build
var dashboard embed.FS

const INDEX = "index.html"

func Server(apiconfig config.Config) {
	var wgServer sync.WaitGroup
	stateError := make(chan error)
	wgServer.Add(1)

	go poolMonitor(apiconfig, stateError)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers, Range, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	fsEmbed := EmbedFolder(dashboard, "dashboard/build", true)
	router.Use(Serve("/", fsEmbed))

	serverRoot, err := fs.Sub(dashboard, "dashboard/build")
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/auth", func(c *gin.Context) {
		err = basicAuth(c, apiconfig)
		if err != nil {
			auth(c)
		} else {
			auth(c)
		}
	})

	router.GET("/api/hello", hello)
	router.GET("/api/monitor/0", AuthMiddleware(), func(ctx *gin.Context) {
		monitor(ctx, apiconfig)
	})

	router.GET("/", func(c *gin.Context) {
		fmt.Printf("URL: %s\n", c.Request.URL.Path)
		c.FileFromFS("/index.html", http.FileSystem(http.FS(serverRoot)))
		c.AbortWithStatus(200)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found",
		})
	})

	/* router.NoRoute(func (c *gin.Context) {
		fmt.Println("%s doesn't exists, redirect on /", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	}) */

	/* 	staticRoot, err := fs.Sub(dashboard, "dashboard/build")
	   	if err != nil {
	   		log.Fatal(err)
	   	} */

	//router.Use(AuthMiddleware())

	router.Run(":9090")

}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "world"})
}

type embedFileSystem struct {
	http.FileSystem
	indexes bool
}

func EmbedFolder(fsEmbed embed.FS, targetPath string, index bool) embedFileSystem {
	fmt.Printf("TargetPath: %s\n", targetPath)
	fmt.Printf("Index: %t\n", index)
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(subFS),
		indexes:    index,
	}
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return false
	}

	// check if indexing is allowed
	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}

func Serve(urlPrefix string, fs embedFileSystem) gin.HandlerFunc {
	fmt.Printf("URL: %s\n", urlPrefix)
	fmt.Printf("FS: %s\n", fs.FileSystem)
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
