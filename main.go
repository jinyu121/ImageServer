package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	port = flag.Int("port", 9420, "Listen port")
)

func main() {
	flag.Parse()

	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	server.Run(fmt.Sprintf("[::]:%d", *port))
}
