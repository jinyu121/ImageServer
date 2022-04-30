package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/util"
)

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run(fmt.Sprintf(":%d", *util.Port))
	if err != nil {
		return
	}

}
