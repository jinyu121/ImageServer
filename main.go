package main

import (
	"embed"
	"flag"
	"fmt"
	"haoyu.love/ImageServer/app/handler/folder_handler"
	"haoyu.love/ImageServer/app/util"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
)

var (
	//go:embed static templates
	assets embed.FS
)

func main() {
	flag.Parse()
	if 0 == flag.NArg() {
		app.Root = "./"
	} else {
		app.Root = flag.Arg(0)
	}
	app.Root, _ = filepath.Abs(app.Root)

	// Ensure the root exists.
	fileInfo, err := os.Stat(app.Root)
	if os.IsNotExist(err) {
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

	var ProcessFn func(*gin.Context)

	if fileInfo.IsDir() {
		ProcessFn = folder_handler.Process
	} else {

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

	_ = appRouter.Run(fmt.Sprintf(":%d", *app.Port))

}
