//go:build linux || darwin || windows
// +build linux darwin windows

package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"haoyu.love/ImageServer/app/handler/folder_handler"
	"haoyu.love/ImageServer/app/handler/lmdb_handler"
	"haoyu.love/ImageServer/app/handler/text_handler"
	"haoyu.love/ImageServer/app/util"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
)

var (
	Version = "Unknown"
	Build   = "Unknown"
)

var (
	//go:embed static templates
	assets embed.FS
)

func InitServer() *gin.Engine {
	// Select proper process function
	var ProcessFn func(*gin.Context)
	if fileInfo, _ := os.Stat(app.Root); fileInfo.IsDir() {
		if ".lmdb" == filepath.Ext(app.Root) {
			lmdb_handler.Init(app.Root)
			ProcessFn = lmdb_handler.Process
		} else {
			ProcessFn = folder_handler.Process
		}
	} else {
		text_handler.Init(app.Root, *app.Column, *app.CustomJsonPath)
		ProcessFn = text_handler.Process
	}

	// Router for the framework itself, such as static files
	frameworkRouter := gin.New()
	frameworkG := frameworkRouter.Group("/_")

	staticFiles, _ := fs.Sub(assets, "static")
	frameworkG.StaticFS("/", http.FS(staticFiles))

	// The general router
	appRouter := gin.Default()
	templateFiles := template.Must(
		template.New("").Funcs(util.TemplateFunction).ParseFS(assets, "templates/*.html"))
	appRouter.SetHTMLTemplate(templateFiles)

	appRouter.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		// Special handling for the
		if strings.HasPrefix(path, "/_/") {
			frameworkRouter.HandleContext(c)
		} else {
			ProcessFn(c)
		}
	})
	return appRouter
}

func main() {
	if "Unknown" != Version {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("ImageServer", Version, "Build", Build)
	app.InitFlag()

	go app.CheckUpdate(Version)

	appRouter := InitServer()

	go func() {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", *app.Port),
			Handler: appRouter,
		}
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
		}
	}()

	listenOn := util.GetIPAddress()
	if len(listenOn) > 0 {
		log.Println("Listening on these addresses:")
		for _, addr := range listenOn {
			if addr.To4() != nil {
				log.Printf("\thttp://%s:%d\n", addr, *app.Port)
			} else {
				log.Printf("\thttp://[%s]:%d\n", addr, *app.Port)
			}
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Bye~")
}
