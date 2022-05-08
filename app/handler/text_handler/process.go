package text_handler

import (
	"bufio"
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
	"github.com/spyzhov/ajson"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

var (
	Data = make([]string, 0)
)

func Init(path string, column int, jsonP string) {
	log.Printf("Scan column %d of file %s", column, path)
	bar := progressbar.Default(-1, "Scanning")

	ext := strings.ToLower(filepath.Ext(path))

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	if ".json" == ext {
		initJson(f, jsonP, bar)
	} else {
		initText(f, ext, column, bar)
	}

	log.Printf("Done! %d records read", len(Data))
}

func initText(f *os.File, ext string, column int, bar *progressbar.ProgressBar) {
	csvReader := csv.NewReader(f)
	if ".csv" != ext {
		csvReader.Comma = '\t'
	}

	for {
		_ = bar.Add(1)
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

func initJson(f *os.File, jsonP string, bar *progressbar.ProgressBar) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_ = bar.Add(1)

		root, err := ajson.Unmarshal(scanner.Bytes())
		if nil != err {
			continue
		}
		nodes, err := root.JSONPath(jsonP)
		if nil != err {
			continue
		}
		for _, node := range nodes {
			s, err := node.GetString()
			if nil != err {
				continue
			}
			Data = append(Data, s)
		}
	}
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
