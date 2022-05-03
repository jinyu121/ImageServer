package lmdb_handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

func Init(path string) {
	log.Printf("Open database %s", path)
	InitDB(path)

	log.Printf("Scan database %s", path)
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
			counter += 1
		}
	})
	log.Printf("Done! %d records read", counter)
}

func Process(c *gin.Context) {
	name := c.Param("path")[1:]
	namePart := strings.Split(name, "/")
	currNode := LmdbTree
	if "" == name {
		processSingleFolder(c)
	} else {
		for ith, k := range namePart {
			if _, ok := currNode.Children[k]; !ok {
				c.String(http.StatusNotFound, "Not found")
				return
			}
			currNode = currNode.Children[k]
			if ith == len(namePart)-1 {
				if currNode.IsFile {
					processFile(c)
				} else {
					processSingleFolder(c)
				}
				return
			}
		}
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

func processSingleFolder(c *gin.Context) {
	name := c.Param("path")[1:]
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	namePart := strings.Split(name, "/")
	currNode := LmdbTree
	if "" != name {
		for _, k := range namePart {
			currNode = currNode.Children[k]
		}
	}

	folderPrev, folderNext, folderParent := "", "", ""
	if nil != currNode.Parent {
		folderPrev, folderNext = GetNeighborFolder(currNode)
		folderParent = GetPath(currNode.Parent)
	}

	folders, files, _ := GetFolderContent(currNode)
	folders, files, pageNum, pageNumMax, pagePrev, pageNext := util.Pagination(*app.PageSize, pageNum, folders, files)

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
