package core

import "strings"

type Breadcrumb struct {
	Show string
	Link string
}
type PageData struct {
	Folders []string
	Files   []string
	Path    []Breadcrumb
}

func NewPageData(folders, files []string, path string) *PageData {
	crumbs := []Breadcrumb{Breadcrumb{Show: "Home", Link: "/"}}
	path = strings.Trim(path, "/")
	if "" != path {
		sps := strings.Split(path, "/")
		for i, item := range sps {
			crumbs = append(crumbs, Breadcrumb{item, "/" + strings.Join(sps[:i+1], "/")})
		}
	}
	return &PageData{Folders: folders, Files: files, Path: crumbs}
}
