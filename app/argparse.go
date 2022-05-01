package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	Port     = flag.Int("Port", 9420, "Listen Port")
	Root     = flag.String("root", "./", "Image folder, or image list file")
	PageSize = flag.Int("page", 1000, "Page size")
	Column   = flag.Int("column", 0, "Column")
)

func Init() {
	flag.Parse()

	*Root, _ = filepath.Abs(*Root)

	if _, err := os.Stat(*Root); os.IsNotExist(err) {
		panic(fmt.Sprintf("Path %s doesn't exists", *Root))
	}

}
