package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app/datasource"
)

type ImageServerHandler struct {
	data datasource.DataSource
}

func NewImageServerHandler(data *datasource.DataSource) *ImageServerHandler {
	handler := &ImageServerHandler{data: *data}
	return handler
}

func (handler *ImageServerHandler) Handle(c *gin.Context) {
	name := c.Param("path")

	fileInfo := handler.data.Stat(name)

	if !fileInfo.Exists {
		c.String(http.StatusNotFound, fmt.Sprintf("Path %s not found", name))
		return
	} else {
		if fileInfo.IsFile {
			handler.processFile(c)
		} else {
			handler.processFolder(c)
		}
	}

}

func (handler *ImageServerHandler) processFile(c *gin.Context) {
	name := c.Param("path")
	content, err := handler.data.GetFile(name)
	if nil != err {
		c.String(http.StatusInternalServerError, "Internal Server Error: %s", err)
	}
	c.Data(http.StatusOK, mimetype.Detect(content).String(), content)
}

func (handler *ImageServerHandler) processFolder(c *gin.Context) {
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	folderNames := []string{c.Param("path")}
	for _, pa := range c.QueryArray("c") {
		stat := handler.data.Stat(pa)
		if !stat.Exists || stat.IsFile {
			continue
		}
		folderNames = append(folderNames, pa)
	}

	contents := make([]datasource.FolderContent, 0)
	for _, name := range folderNames {
		content, err := handler.data.GetFolder(name)
		if nil != err {
			c.String(http.StatusInternalServerError, "Internal Server Error: %s", err)
		}
		contents = append(contents, content)
	}

	aligned := DeduplicateFolderContent(&contents)
	pagination := Paginate(&aligned, &contents, *PageSize, pageNum, GetCurrentUrl(c))

	navigation := datasource.Navigation{}
	if len(contents) == 1 {
		content := contents[0]
		navigation = *handler.data.GetNeighbor(contents[0].Name)

		c.HTML(http.StatusOK, "list.html", gin.H{
			"content":    content,
			"pagination": pagination,
			"navigation": navigation,
		})
	} else {
		c.HTML(http.StatusOK, "compare.html", gin.H{
			"contents":    contents,
			"pagination":  pagination,
			"navigation":  navigation,
			"aligned":     aligned,
			"columnWidth": 90. / len(contents), // Since the label column will take 10% of width
		})
	}
}
