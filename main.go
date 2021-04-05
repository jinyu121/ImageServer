package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/jinyu121/ImageServer/app"
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
	TemplateList    = app.GetTemplate(TemplateFiles, "static/template/list.tmpl")
	TemplateCompare = app.GetTemplate(TemplateFiles, "static/template/compare.tmpl")
	Template404     = app.GetTemplate(TemplateFiles, "static/template/404.tmpl")
)

func processDirectory(w http.ResponseWriter, req *http.Request) {
	rootInfo, err := os.Stat(path.Join(*root, strings.TrimRight(req.URL.Path, "/")))
	if nil != err {
		process404(w, req)
		return
	}

	if rootInfo.Mode().IsDir() {
		comparePaths, ok := req.URL.Query()["c"]
		if ok {
			paths := []string{req.URL.Path}
			paths = append(paths, comparePaths...)
			multiPageData := app.MakeMultiPageData(paths, *root)
			TemplateCompare.Execute(w, multiPageData)
		} else {
			pageData := app.MakePageData(req.URL.Path, *root)
			TemplateList.Execute(w, pageData)
		}
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

	fmt.Println("Server is ready to handle requests at port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))
}
