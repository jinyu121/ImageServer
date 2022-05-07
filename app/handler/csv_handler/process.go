package csv_handler

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/schollz/progressbar/v3"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

var (
	Data = make([]string, 0)
)

func Init(path string, column int) {
	log.Printf("Scan column %d of file %s", column, path)
	bar := progressbar.Default(-1, "Scanning")

	ext := strings.ToLower(filepath.Ext(path))

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	csvReader := csv.NewReader(f)
	if ".csv" != ext {
		csvReader.Comma = '\t'
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(rec) < column {
			continue
		}
		Data = append(Data, rec[column])
		_ = bar.Add(1)
	}
	log.Printf("Done! %d records read", len(Data))
}

func Process(c *gin.Context) {
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	contents := []util.FolderContent{{
		Name:    "",
		Folders: make([]string, 0),
		Files:   Data,
	},
	}
	pagination := util.Paginate(&contents, *app.PageSize, pageNum, util.GetCurrentUrl(c))
	navigation := util.Navigation{}

	content := contents[0]
	c.HTML(http.StatusOK, "list.html", gin.H{
		"content":    content,
		"pagination": pagination,
		"navigation": navigation,
	})
}
