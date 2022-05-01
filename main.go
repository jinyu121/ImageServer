package main

import (
	"embed"
	"fmt"
	"haoyu.love/ImageServer/app/handler/folder_handler"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
)

var (
	//go:embed static/css static/images static/js
	StaticFiles embed.FS
)

func main() {
	app.Init()

	var ProcessFn func(*gin.Context)
	fileInfo, _ := os.Stat(*app.Root)
	if fileInfo.IsDir() {
		ProcessFn = folder_handler.Process
	} else {

	}

	// Router for the framework itself, such as static files
	frameworkRouter := gin.New()
	frameworkG := frameworkRouter.Group("/_")
	frameworkG.StaticFS("/static", http.FS(StaticFiles))

	// The general router
	appRouter := gin.Default()
	appRouter.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		// Special handling for the
		if strings.HasPrefix(path, "/_/") {
			frameworkRouter.HandleContext(c)
		} else {
			ProcessFn(c)
		}
	})

	_ = appRouter.Run(fmt.Sprintf(":%d", *app.Port))

}
