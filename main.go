package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"haoyu.love/ImageServer/app/handler/csv_handler"
	"haoyu.love/ImageServer/app/handler/folder_handler"
	"haoyu.love/ImageServer/app/handler/lmdb_handler"
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

func InitFlag() {
	if "Unknown" != Version {
		gin.SetMode(gin.ReleaseMode)
	}
	flag.Parse()
	if 0 == flag.NArg() {
		app.Root = "./"
	} else {
		app.Root = flag.Arg(0)
	}
	app.Root, _ = filepath.Abs(app.Root)

	if *app.Column < 0 {
		*app.Column = 0
	}

	// Ensure the root exists.
	if _, err := os.Stat(app.Root); os.IsNotExist(err) {
		panic(fmt.Sprintf("Path %s doesn't exists", app.Root))
	}

	// Process Extension
	if "" != *app.ExtCustom {
		ext := strings.Split(strings.ToLower(*app.ExtCustom), ",")
		util.ArrayToSet(app.FileExtension, ext)
	} else {
		util.ArrayToSet(app.FileExtension, app.DefaultImageExt)
		util.ArrayToSet(app.FileExtension, app.DefaultAudioExt)
		util.ArrayToSet(app.FileExtension, app.DefaultVideoExt)
	}
}

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
		csv_handler.Init(app.Root, *app.Column)
		ProcessFn = csv_handler.Process
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
	log.Println("ImageServer", Version, "Build", Build)
	InitFlag()
	appRouter := InitServer()
	_ = appRouter.Run(fmt.Sprintf(":%d", *app.Port))
}
