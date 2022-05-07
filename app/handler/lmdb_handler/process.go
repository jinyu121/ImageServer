package lmdb_handler

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/schollz/progressbar/v3"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

func Init(path string) {
	log.Printf("Open database %s", path)
	InitDB(path)

	log.Printf("Scan database %s", path)
	bar := progressbar.Default(-1, "Scanning")
	counter := 0
	_ = LmdbEnv.View(func(txn *lmdb.Txn) (err error) {
		cur, err := txn.OpenCursor(LmdbDBI)
		if err != nil {
			return err
		}
		defer cur.Close()

		for {
			k, _, err := cur.Get(nil, nil, lmdb.Next)
			if lmdb.IsNotFound(err) {
				return nil
			}
			if err != nil {
				return err
			}

			AddToTree(string(k))
			_ = bar.Add(1)
			counter += 1
		}
	})
	log.Printf("Scan Done! %d records read", counter)
}

func Process(c *gin.Context) {
	name := c.Param("path")
	currNode, err := GetNode(name)
	if nil != err {
		c.String(http.StatusNotFound, "Not found")
	} else if currNode.IsFile {
		processFile(c)
	} else {
		processFolder(c)
	}
}

func processFile(c *gin.Context) {
	name := c.Param("path")[1:]

	_ = LmdbEnv.View(func(txn *lmdb.Txn) (err error) {
		v, err := txn.Get(LmdbDBI, []byte(name))
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: %s", err)
			return nil
		}
		c.Data(http.StatusOK, mimetype.Detect(v).String(), v)
		return nil
	})
}

func processFolder(c *gin.Context) {
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	folderNames := []string{c.Param("path")}
	for _, fdr := range c.QueryArray("c") {
		folderNames = append(folderNames, fdr)
	}
	contents := make([]util.FolderContent, 0)
	for _, name := range folderNames {
		node, err := GetNode(name)
		if nil != err {
			continue
		}
		content, err := GetFolderContent(node)
		if nil != err {
			continue
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
		currNode, _ := GetNode(content.Name)

		if nil != currNode.Parent {
			navigation.Prev, navigation.Next = GetNeighborFolder(currNode)
			navigation.Current = filepath.Dir(navigation.Current)
		}

		c.HTML(http.StatusOK, "list.html", gin.H{
			"content":    content,
			"pagination": pagination,
			"navigation": navigation,
		})
	} else {
		c.HTML(http.StatusOK, "compare.html", gin.H{
			"contents":   contents,
			"pagination": pagination,
			"navigation": navigation,
			"aligned":    aligned,
		})
	}
}
