package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"haoyu.love/ImageServer/app/util"
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
	FileExtension = make(map[string]struct{})
)

func InitFlag() {
	flag.Parse()
	if 0 == flag.NArg() {
		Root = "./"
	} else {
		Root = flag.Arg(0)
	}
	Root, _ = filepath.Abs(Root)

	if *Column < 0 {
		*Column = 0
	}

	// Ensure the root exists.
	if _, err := os.Stat(Root); os.IsNotExist(err) {
		panic(fmt.Sprintf("Path %s doesn't exists", Root))
	}

	// Process Extension
	if "" != *CustomExt {
		if "*" != *CustomExt {
			ext := strings.Split(strings.ToLower(*CustomExt), ",")
			util.ArrayToSet(FileExtension, ext)
		}
	} else {
		util.ArrayToSet(FileExtension, util.DefaultImageExt)
		util.ArrayToSet(FileExtension, util.DefaultAudioExt)
		util.ArrayToSet(FileExtension, util.DefaultVideoExt)
	}
}
