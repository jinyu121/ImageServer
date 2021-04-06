package app

import (
	"embed"
	"html/template"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type PageData struct {
	Path       string
	Folders    []string
	Images     []string
	Videos     []string
	Files      []string
	PreFolder  string
	NextFolder string
}

func MakePageData(urlPath, root string) *PageData {
	realPath := strings.TrimSuffix(root, "/") + urlPath
	folders, files, _ := GetFolderContent(realPath)
	images := FilterFiles(files, IsImageFile)
	videos := FilterFiles(files, IsVideoFile)
	preFolder, nextFolder := GetNeighborFolder(realPath, root, 1)
	pageData := PageData{
		Path:       urlPath,
		Folders:    RemoveLeft(folders, root),
		Images:     RemoveLeft(images, root),
		Videos:     RemoveLeft(videos, root),
		PreFolder:  preFolder,
		NextFolder: nextFolder,
	}
	return &pageData
}

type MultiPageData struct {
	Pages []PageData
}

func MakeMultiPageData(urlPaths []string, root string) MultiPageData {
	pd := make([]PageData, 0)
	for _, urlPath := range urlPaths {
		pd = append(pd, *MakePageData(urlPath, root))
	}
	return MultiPageData{Pages: pd}
}

func (mpd MultiPageData) Folders() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Folders
	}
	return stringAlign(tmp)
}

func (mpd MultiPageData) Images() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Images
	}
	return stringAlign(tmp)
}

func (mpd MultiPageData) Videos() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Videos
	}
	return stringAlign(tmp)
}

func (mpd MultiPageData) Files() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Files
	}
	return stringAlign(tmp)
}

func stringAlign(data [][]string) [][]string {
	// How many arrays
	n := len(data)
	// Golang doesn't have dataset of set, so we have to use map
	filesSet := make([]map[string]string, n)
	// fileSet stores all the files
	fileSet := make(map[string]bool)
	// Record each array
	for i, dataList := range data {
		filesSet[i] = make(map[string]string)
		for _, f := range dataList {
			name := filepath.Base(f)
			fileSet[name] = true
			filesSet[i][name] = f
		}
	}
	// Now we can get a non-duplicated file list
	var i = 0
	fileList := make([]string, len(fileSet))
	for k := range fileSet {
		fileList[i] = k
		i++
	}
	sort.Strings(fileList)

	// Make final result
	result := make([][]string, len(fileSet))
	for i, k := range fileList {
		line := make([]string, n+1)
		line[0] = k
		for j, fileSetItem := range filesSet {
			v, ok := fileSetItem[k]
			if ok {
				line[j+1] = v
			} else {
				line[j+1] = ""
			}
		}
		result[i] = line
	}

	return result
}

func GetTemplate(storage embed.FS, fileList ...string) *template.Template {
	var tpl *template.Template
	fileList = append(fileList, "static/template/base.tmpl")
	tpl, _ = template.New(path.Base(fileList[0])).Funcs(templateFunction).ParseFS(storage, fileList...)
	return tpl
}

var templateFunction = template.FuncMap{
	"pathToName": func(p string) string {
		return path.Base(p)
	},
	"lastOne": func(arr []interface{}) interface{} {
		if 0 == len(arr) {
			return nil
		}
		return arr[len(arr)-1]
	},
	"breadCrumb": func(root string) []string {
		crumb := make([]string, 0)
		root = strings.Trim(root, "/")
		if "" != root {
			sps := strings.Split(root, "/")
			for i := range sps {
				crumb = append(crumb, "/"+strings.Join(sps[:i+1], "/"))
			}
		}
		return crumb
	},
}
