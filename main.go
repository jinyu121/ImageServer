package main

import (
	"embed"
	"flag"
	"fmt"
	"haoyu.love/ImageServer/app"
	"net/http"
	"os"
	"path/filepath"
)

const (
	CODENAME string = "ImageServer"
)

var (
	Version     string = "<undefined>"
	Build       string = "<undefined>"
	VersionLong string = fmt.Sprintf("%s %s (Build %s)", CODENAME, Version, Build)
)

var (
	//go:embed static/css static/images static/js
	StaticFiles embed.FS
	//go:embed static/template
	TemplateFiles embed.FS

	// Args
	root         = flag.String("root", "./", "Image folder, or image list file")
	port         = flag.Int("port", 9420, "Listen port")
	printVersion = flag.Bool("version", false, "Print version and exit")

	// Global variables
	fileLines  []string
	fileServer func(w http.ResponseWriter, req *http.Request) = nil
)

func main() {
	flag.Parse()

	fmt.Println(VersionLong)

	if *printVersion {
		os.Exit(0)
	}

	*root, _ = filepath.Abs(*root)
	_, err := os.Stat(*root)
	if nil != err {
		fmt.Print("Path not exist")
		return
	}

	server := app.NewServer(*root, TemplateFiles, StaticFiles)

	fmt.Println("Server is ready to handle requests at port", *port)
	server.Run(*port)

}
