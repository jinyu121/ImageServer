package app

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app/datasource"
	"haoyu.love/ImageServer/app/filter"
)

var (
	Root     = "./"
	Port     = flag.Int("port", 9420, "Listen Port")
	PageSize = flag.Int("page", 1000, "Page size")
	Column   = flag.Int("column", 0, "Column")
	Filter   = flag.String("filter", filter.PredefineDefault, "Filter")
)

func InitFlag() {
	flag.Parse()
	if flag.NArg() == 0 {
		Root = "./"
	} else {
		Root = flag.Arg(0)
	}
	Root, _ = filepath.Abs(Root)

	if *Column < 0 {
		*Column = 0
	}

	// Ensure the root exists.
	if _, err := os.Stat(Root); os.IsNotExist(err) {
		panic(fmt.Sprintf("Path %s doesn't exists", Root))
	}

}

func InitServer(assets embed.FS) *gin.Engine {
	// Select proper data source
	var data datasource.DataSource
	var flt filter.Filter

	if fileInfo, _ := os.Stat(Root); fileInfo.IsDir() {
		if filepath.Ext(Root) == ".lmdb" {
			flt = filter.NewNoFilter()
			data = datasource.NewLmdbDataSource(Root, &flt)
		} else {
			flt = filter.NewFileExtFilter(*Filter)
			data = datasource.NewFolderDataSource(Root, &flt)
		}
	} else {
		flt = filter.NewJsonFilter(*Filter)
		data = datasource.NewTextFileDataSource(Root, &flt, *Column)
	}

	handler := NewImageServerHandler(&data)

	// Router for the framework itself, such as static files
	frameworkRouter := gin.New()
	frameworkG := frameworkRouter.Group("/_")

	staticFiles, _ := fs.Sub(assets, "static")
	frameworkG.StaticFS("/", http.FS(staticFiles))

	// The general router
	appRouter := gin.Default()
	templateFiles := template.Must(
		template.New("").Funcs(TemplateFunction).ParseFS(assets, "templates/*.html"))
	appRouter.SetHTMLTemplate(templateFiles)

	appRouter.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		// Special handling for the static files
		if strings.HasPrefix(path, "/_/") {
			frameworkRouter.HandleContext(c)
		} else {
			handler.Handle(c)
		}
	})
	return appRouter
}
