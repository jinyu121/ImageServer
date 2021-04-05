package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/jinyu121/ImageServer/app"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	//go:embed static/css static/images static/js
	StaticFiles embed.FS
	//go:embed static/template
	TemplateFiles embed.FS

	// Args
	root = flag.String("root", "./", "Image folder, or image list file")
	port = flag.Int("port", 9420, "Listen port")

	// Global variables
	fileLines  []string
	fileServer func(w http.ResponseWriter, req *http.Request) = nil
)

var (
	TemplateList = parseTemplate("static/template/list.tmpl")
	Template404  = parseTemplate("static/template/404.tmpl")
)

func parseTemplate(fileList ...string) *template.Template {
	var tpl *template.Template
	fileList = append(fileList, "static/template/base.tmpl")
	tpl, _ = template.New(path.Base(fileList[0])).Funcs(app.TemplateFunction).ParseFS(TemplateFiles, fileList...)
	return tpl
}

func processDirectory(w http.ResponseWriter, req *http.Request) {
	realPath := path.Join(*root, strings.Trim(req.URL.Path, "/"))

	rootInfo, err := os.Stat(realPath)
	if nil != err {
		process404(w, req)
		return
	}

	if rootInfo.Mode().IsDir() {
		folders, files, _ := app.GetFolderContent(realPath)
		images := app.FilterFiles(files, app.IsImageFile)
		videos := app.FilterFiles(files, app.IsVideoFile)
		pageData := app.PageData{
			Path:    req.URL.Path,
			Folders: app.RemoveLeft(folders, *root),
			Images:  app.RemoveLeft(images, *root),
			Videos:  app.RemoveLeft(videos, *root),
		}
		TemplateList.Execute(w, pageData)
	} else {
		fileServer(w, req)
	}
}

func processFile(w http.ResponseWriter, req *http.Request) {
	pageData := app.PageData{Path: "/", Files: fileLines}
	TemplateList.Execute(w, pageData)
}

func process404(w http.ResponseWriter, req *http.Request) {
	pageData := app.PageData{Path: req.URL.Path}
	Template404.Execute(w, pageData)
}

func main() {
	flag.Parse()

	rootInfo, err := os.Stat(*root)
	if nil != err {
		return
	}
	if rootInfo.Mode().IsDir() {
		if !strings.HasPrefix(*root, "/") {
			*root = *root + "/"
		}
		fileServer = http.FileServer(http.Dir(*root)).ServeHTTP
		http.HandleFunc("/", processDirectory)
	} else {
		fileLines, _ = app.GetTextContent(*root)
		http.HandleFunc("/", processFile)
	}

	staticFilesClean, _ := fs.Sub(StaticFiles, "static")
	http.HandleFunc("/_/", http.StripPrefix("/_/", http.FileServer(http.FS(staticFilesClean))).ServeHTTP)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))
}
