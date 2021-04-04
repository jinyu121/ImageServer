package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/jinyu121/ImageServer/core"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	//go:embed static/css static/image static/js
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

func processDirectory(w http.ResponseWriter, req *http.Request) {
	realPath := path.Join(*root, strings.Trim(req.URL.Path, "/"))

	rootInfo, err := os.Stat(realPath)
	if nil != err {
		process404(w, req)
		return
	}

	if rootInfo.Mode().IsDir() {
		folders, files, _ := core.GetFolderContent(realPath)
		files = core.FilterImageFiles(files)

		pageData := core.NewPageData(core.RemoveLeft(folders, *root), core.RemoveLeft(files, *root), req.URL.Path)

		t, _ := template.ParseFS(TemplateFiles, "static/template/list.tmpl", "static/template/base.tmpl")
		t.Execute(w, pageData)
	} else {
		fileServer(w, req)
	}
}

func processFile(w http.ResponseWriter, req *http.Request) {
	pageData := core.NewPageData(make([]string, 0), fileLines, "/")

	t, _ := template.ParseFS(TemplateFiles, "static/template/list.tmpl", "static/template/base.tmpl")
	t.Execute(w, pageData)
}

func process404(w http.ResponseWriter, req *http.Request) {
	pageData := core.NewPageData(make([]string, 0), make([]string, 0), req.URL.Path)
	t, _ := template.ParseFS(TemplateFiles, "static/template/404.tmpl", "static/template/base.tmpl")
	t.Execute(w, pageData)
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
		fileLines, _ = core.GetTextContent(*root)
		http.HandleFunc("/", processFile)
	}

	staticFilesClean, _ := fs.Sub(StaticFiles, "static")
	http.HandleFunc("/_/", http.StripPrefix("/_/", http.FileServer(http.FS(staticFilesClean))).ServeHTTP)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))
}
