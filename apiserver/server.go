package apiserver

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/aescanero/openldap-node/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

//go:embed dashboard/build
var dashboard embed.FS

func Server(apiconfig config.Config) {
	serverRoot, err := fs.Sub(dashboard, "dashboard/build")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	router.GET("/auth", func(c *gin.Context) {
		err = basicAuth(c, apiconfig)
		if err != nil {
			auth(c)
		}
	})
	router.Use(AuthMiddleware())

	router.StaticFS("/dashboard", http.FileSystem(http.FS(serverRoot)))
	router.GET("/", AuthMiddleware(), func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/auth"
		router.HandleContext(ctx)
	})

	router.GET("/api/hello", hello)
	router.GET("/api/monitor", AuthMiddleware(), func(ctx *gin.Context) {
		err = basicAuth(ctx, apiconfig)
		if err != nil {
			monitor(ctx, apiconfig)
		}
	})
	router.Run(":9090")
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "world"})
}

func auth(ctx *gin.Context) {
	username := ctx.PostForm("username")
	//password := c.PostForm("password")
	claims := jwt.MapClaims{
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("world"))
	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token not found"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("world"), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ProtectedHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Protected endpoint accessed"})
}
