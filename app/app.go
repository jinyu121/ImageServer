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
)

var (
	Root           = "./"
	Port           = flag.Int("port", 9420, "Listen Port")
	PageSize       = flag.Int("page", 1000, "Page size")
	Column         = flag.Int("column", 0, "Column")
	CustomExt      = flag.String("ext", "", "File extensions")
	CustomJsonPath = flag.String("json", "", "JsonPath if you are using json file")
)

var (
	FileExtension = make(map[string]struct{})
)

func InitFlag() {
	flag.Parse()
	if 0 == flag.NArg() {
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

	// Process Extension
	if "" != *CustomExt {
		if "*" != *CustomExt {
			ext := strings.Split(strings.ToLower(*CustomExt), ",")
			ArrayToSet(FileExtension, ext)
		}
	} else {
		ArrayToSet(FileExtension, DefaultImageExt)
		ArrayToSet(FileExtension, DefaultAudioExt)
		ArrayToSet(FileExtension, DefaultVideoExt)
	}
}

func InitServer(assets embed.FS) *gin.Engine {
	// Select proper data source
	var data DataSource

	if fileInfo, _ := os.Stat(Root); fileInfo.IsDir() {
		if ".lmdb" == filepath.Ext(Root) {
			data = datasource.NewLmdbDataSource(Root)
		} else {
			data = datasource.NewFolderDataSource(Root)
		}
	} else {
		data = datasource.NewTextFileDataSource(Root, *CustomJsonPath, *Column)
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
		// Special handling for the
		if strings.HasPrefix(path, "/_/") {
			frameworkRouter.HandleContext(c)
		} else {
			handler.Handle(c)
		}
	})
	return appRouter
}
