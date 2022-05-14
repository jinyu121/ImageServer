package folder_handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

func Process(c *gin.Context) {
	name := c.Param("path")
	fullPath := filepath.Join(app.Root, name)
	fullPath, _ = filepath.Abs(fullPath)

	if ok, _ := IsSubFolder(app.Root, fullPath); !ok {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Not found")
	} else if fileInfo.IsDir() {
		processFolder(c)
	} else {
		processFile(c)
	}
}

func processFile(c *gin.Context) {
	name := filepath.Join(app.Root, c.Param("path"))
	c.File(name)
}

func processFolder(c *gin.Context) {
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	folderNames := []string{filepath.Join(app.Root, c.Param("path"))}
	for _, fdr := range c.QueryArray("c") {
		folderNames = append(folderNames, filepath.Join(app.Root, fdr))
	}

	contents := make([]util.FolderContent, 0)
	for _, name := range folderNames {
		content, err := GetFolderContent(name)
		if nil != err {
			if os.IsNotExist(err) {
				c.String(http.StatusNotFound, fmt.Sprintf("Path %s not found", name))
			} else {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Error while reading folder %s", name))
			}
			return
		}
		content.FilterTargetFile()
		contents = append(contents, content)
	}
	aligned := util.AlignContent(&contents)
	pagination := util.Paginate(&contents, *app.PageSize, pageNum, util.GetCurrentUrl(c))
	aligned = (*util.Paginate(&[]util.FolderContent{aligned}, *app.PageSize, pageNum, "").Content)[0]

	navigation := util.Navigation{}
	if 1 == len(contents) {
		content := contents[0]
		navigation.Current = content.Name
		if contents[0].Name != app.Root {
			navigation.Prev, navigation.Next = GetNeighborFolder(contents[0].Name)
			navigation.Parent = filepath.Dir(navigation.Current)
		}
		navigation.RemovePrefix(app.Root)

		content.RemovePrefix(app.Root)
		c.HTML(http.StatusOK, "list.html", gin.H{
			"content":    content,
			"pagination": pagination,
			"navigation": navigation,
		})
	} else {
		for i := range contents {
			contents[i].RemovePrefix(app.Root)
		}
		c.HTML(http.StatusOK, "compare.html", gin.H{
			"contents":   contents,
			"pagination": pagination,
			"navigation": navigation,
			"aligned":    aligned,
		})
	}
}
