package app

import (
	"flag"
)

var (
	Version = "Unknown"
	Build   = "Unknown"
)

var (
	Root           = "./"
	Port           = flag.Int("port", 9420, "Listen Port")
	PageSize       = flag.Int("page", 1000, "Page size")
	Column         = flag.Int("column", 0, "Column")
	CustomExt      = flag.String("ext", "", "File extensions")
	CustomJsonPath = flag.String("json", "", "JsonPath if you are using json file")
)

var (
	DefaultImageExt = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".svg", ".webp", ".ico"}
	DefaultVideoExt = []string{".mp4", ".mkv", ".mov", ".wmv", ".flv", ".avi", ".rmvb", ".mpg", ".mpeg", ".m4v", ".3gp", ".3g2"}
	DefaultAudioExt = []string{".mp3", ".wav", ".wma", ".ogg", ".flac"}
)

var (
	FileExtension = make(map[string]struct{})
)
