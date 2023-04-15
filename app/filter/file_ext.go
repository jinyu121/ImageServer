package filter

import (
	"path/filepath"
	"strings"
)

var (
	DefaultAudioExt = []string{".mp3", ".wav", ".wma", ".ogg", ".flac"}
	DefaultImageExt = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".svg", ".webp", ".ico"}
	DefaultVideoExt = []string{".mp4", ".mkv", ".mov", ".wmv", ".flv", ".avi", ".rmvb", ".mpg", ".mpeg", ".m4v", ".3gp", ".3g2"}
)

type FileExtFilter struct {
	ext map[string]interface{}
}

func NewFileExtFilter(ext string) *FileExtFilter {
	filter := &FileExtFilter{ext: make(map[string]interface{})}
	tmp := make([]string, 0)
	if PredefineDefault == ext || PredefineImageExt == ext {
		tmp = DefaultImageExt
	} else if PredefineVideoExt == ext {
		tmp = DefaultVideoExt
	} else if PredefineAudioExt == ext {
		tmp = DefaultAudioExt
	} else {
		tmp = strings.Split(strings.ToLower(ext), ",")
	}

	for _, v := range tmp {
		filter.ext[v] = struct{}{}
	}

	return filter
}

func (f *FileExtFilter) Filter(fileName string) bool {
	if 0 == len(f.ext) {
		return true
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	_, ok := f.ext[ext]
	return ok
}

func (f *FileExtFilter) Extract(line string) []string {
	return []string{line}
}
