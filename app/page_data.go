package app

import "strings"

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
	return AlignStringArrays(tmp)
}

func (mpd MultiPageData) Images() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Images
	}
	return AlignStringArrays(tmp)
}

func (mpd MultiPageData) Videos() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Videos
	}
	return AlignStringArrays(tmp)
}

func (mpd MultiPageData) Files() [][]string {
	tmp := make([][]string, len(mpd.Pages))
	for i, v := range mpd.Pages {
		tmp[i] = v.Files
	}
	return AlignStringArrays(tmp)
}
