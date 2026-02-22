package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-minstack/core"
	mgin "github.com/go-minstack/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from MinStack!"})
	})
}

func main() {
	app := core.New(mgin.Module())
	app.Invoke(registerRoutes)
	app.Run()
}
