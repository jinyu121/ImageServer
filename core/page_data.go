package core

import "strings"

type PageData struct {
	Path    string
	Folders []string
	Files   []string
}

func NewPageData(folders []string, files []string, root string) *PageData {
	crumbs := make([]string, 0)
	root = strings.Trim(root, "/")
	if "" != root {
		sps := strings.Split(root, "/")
		for i := range sps {
			crumbs = append(crumbs, "/"+strings.Join(sps[:i+1], "/"))
		}
	}
	return &PageData{Path: root, Folders: folders, Files: files}
}
