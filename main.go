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

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.Run(fmt.Sprintf("[::]:%d", *port))
}
