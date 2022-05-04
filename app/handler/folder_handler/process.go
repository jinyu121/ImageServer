package folder_handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

func Process(c *gin.Context) {
	name := c.Param("path")
	fullPath := filepath.Join(app.Root, name)
	fullPath, _ = filepath.Abs(fullPath)

	if ok, _ := util.IsSubFolder(app.Root, fullPath); !ok {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Not found")
	} else if fileInfo.IsDir() {
		processSingleFolder(c)
	} else {
		processFile(c)
	}
}

func processFile(c *gin.Context) {
	name := filepath.Join(app.Root, c.Param("path"))
	c.File(name)
}

func processSingleFolder(c *gin.Context) {
	name := filepath.Join(app.Root, c.Param("path"))
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	folders, files, err := GetFolderContent(name)
	files = FilterTargetFile(files)

	folders, files, pageNum, pageNumMax, pagePrev, pageNext := util.Pagination(*app.PageSize, pageNum, folders, files)

	folders = util.RemoveLeft(app.Root, folders)
	files = util.RemoveLeft(app.Root, files)

	folderPrev, folderNext, folderParent := "", "", ""
	if name != app.Root {
		folderPrev, folderNext = GetNeighborFolder(name)
		folderPrev = strings.TrimPrefix(folderPrev, app.Root)
		folderNext = strings.TrimPrefix(folderNext, app.Root)

		folderParent = filepath.Dir(name)
		folderParent = strings.TrimPrefix(folderParent, app.Root)
		if "" == folderParent {
			folderParent = "/"
		}
	}

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
