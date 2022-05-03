package csv_handler

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	Data = make([]string, 0)
)

func Init(path string, column int) {
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
	}
}

func Process(c *gin.Context) {
	folders := make([]string, 0)
	files := Data

	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}
	folders, files, pageNum, pageNumMax, pagePrev, pageNext := util.Pagination(*app.PageSize, pageNum, folders, files)

	folderPrev, folderNext, folderParent := "", "", ""

	navigation := gin.H{
		"path":   c.Param("path"),
		"prev":   folderPrev,
		"next":   folderNext,
		"parent": folderParent,
	}

	pagination := gin.H{
		"num":  pageNum,
		"max":  pageNumMax,
		"prev": pagePrev,
		"next": pageNext,
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"folders":    folders,
		"files":      files,
		"pagination": pagination,
		"navigation": navigation,
	})
}
