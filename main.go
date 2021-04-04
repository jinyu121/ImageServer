package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	//go:embed web
	fs embed.FS

	port = flag.Int("port", 9420, "Listen port")
)

func index(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFS(fs, "web/template/index.tmpl", "web/template/base.tmpl")
	t.Execute(w, nil)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))
}
