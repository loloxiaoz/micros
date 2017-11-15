package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Handle all requests using net/http
	http.Handle("/", r)
	r.Run()
}
