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

	folders := make([]string, 0)
	files := Data

	folders, files, pageNum, pageNumMax, pagePrev, pageNext := util.Pagination(*app.PageSize, pageNum, folders, files)

	folderPrev, folderNext, folderParent := "", "", ""

	c.HTML(http.StatusOK, "list.html", gin.H{
		"folders": folders,
		"files":   files,
		"pagination": gin.H{
			"num":  pageNum,
			"max":  pageNumMax,
			"prev": pagePrev,
			"next": pageNext,
		},
		"navigation": gin.H{
			"path":   c.Param("path"),
			"prev":   folderPrev,
			"next":   folderNext,
			"parent": folderParent,
		},
	})
}
