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
	folderNames = append(folderNames, c.QueryArray("c")...)

	contents := make([]datasource.FolderContent, 0)
	for _, name := range folderNames {
		content, err := handler.data.GetFolder(name)
		if nil != err {
			c.String(http.StatusInternalServerError, "Internal Server Error: %s", err)
		}
		content.FilterTargetFile(FileExtension)
		contents = append(contents, content)
	}
	aligned := AlignContent(&contents)
	pagination := Paginate(&contents, *PageSize, pageNum, GetCurrentUrl(c))
	aligned = (*Paginate(&[]datasource.FolderContent{aligned}, *PageSize, pageNum, "").Content)[0]

	navigation := datasource.Navigation{}
	if 1 == len(contents) {
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
			"columnWidth": 90. / len(contents),
		})
	}
}
