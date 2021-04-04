package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

var (
	//go:embed web/image web/static
	staticFiles embed.FS
	//go:embed web/template
	templateFiles embed.FS
	port          = flag.Int("port", 9420, "Listen port")
)

func index(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFS(templateFiles, "web/template/index.tmpl", "web/template/base.tmpl")
	t.Execute(w, nil)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", index)

	staticFilesClean, _ := fs.Sub(staticFiles, "web")
	http.HandleFunc("/_/", http.StripPrefix("/_/", http.FileServer(http.FS(staticFilesClean))).ServeHTTP)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))
}
