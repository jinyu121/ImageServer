package app

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Server struct {
	templates      map[string]*template.Template
	templateStatic fs.FS
	root           string
	fileServer     func(http.ResponseWriter, *http.Request)
	lines          []string
}

func parseTemplate(templateFile embed.FS) map[string]*template.Template {
	templates := make(map[string]*template.Template)
	templates["list"] = GetTemplate(templateFile, "static/template/list.tmpl")
	templates["compare"] = GetTemplate(templateFile, "static/template/compare.tmpl")
	templates["404"] = GetTemplate(templateFile, "static/template/404.tmpl")
	return templates
}

func NewServer(root string, templateFile embed.FS, templateStatic embed.FS) *Server {
	server := Server{
		templates: parseTemplate(templateFile),
		root:      root,
	}

	templateStaticClean, _ := fs.Sub(templateStatic, "static")
	server.templateStatic = templateStaticClean
	http.HandleFunc("/_/", http.StripPrefix("/_/", http.FileServer(http.FS(server.templateStatic))).ServeHTTP)

	rootInfo, _ := os.Stat(root)
	if rootInfo.Mode().IsDir() {
		if !strings.HasPrefix(root, "/") {
			root = root + "/"
		}
		server.fileServer = http.FileServer(http.Dir(root)).ServeHTTP
		http.HandleFunc("/", server.process)
	} else {
		server.lines, _ = GetTextContent(root)
		http.HandleFunc("/", server.processFile)
	}

	return &server
}

func (server *Server) Run(port int) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", port), nil))
}

func (server *Server) process(w http.ResponseWriter, req *http.Request) {
	rootInfo, err := os.Stat(path.Join(server.root, strings.TrimRight(req.URL.Path, "/")))
	if nil != err {
		server.process404(w, req)
		return
	}

	if rootInfo.Mode().IsDir() {
		comparePaths, ok := req.URL.Query()["c"]
		if ok {
			paths := []string{req.URL.Path}
			paths = append(paths, comparePaths...)
			multiPageData := MakeMultiPageData(paths, server.root)
			server.templates["compare"].Execute(w, multiPageData)
		} else {
			pageData := MakePageData(req.URL.Path, server.root)
			server.templates["list"].Execute(w, pageData)
		}
	} else {
		server.fileServer(w, req)
	}
}

func (server *Server) processFile(w http.ResponseWriter, req *http.Request) {
	pageData := PageData{Path: "/", Images: server.lines}
	server.templates["list"].Execute(w, pageData)
}

func (server *Server) process404(w http.ResponseWriter, req *http.Request) {
	pageData := PageData{Path: req.URL.Path}
	server.templates["404"].Execute(w, pageData)
}
