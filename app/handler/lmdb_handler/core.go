package lmdb_handler

import (
	"path"
	"sort"
	"strings"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"haoyu.love/ImageServer/app/util"
)

type Node struct {
	Name     string
	IsFile   bool
	Parent   *Node
	Children map[string]*Node
}

var (
	LmdbEnv  *lmdb.Env
	LmdbDBI  lmdb.DBI
	LmdbTree = &Node{Name: "/", IsFile: false, Parent: nil, Children: make(map[string]*Node)}
)

func InitDB(path string) {
	LmdbEnv, _ = lmdb.NewEnv()
	_ = LmdbEnv.SetMaxDBs(1)
	_ = LmdbEnv.SetMapSize(1 << 42)
	err := LmdbEnv.Open(path, 0, 644)

	if nil != err {
		panic(err)
	}
	err = LmdbEnv.Update(func(txn *lmdb.Txn) (err error) {
		LmdbDBI, err = txn.OpenRoot(lmdb.Create)
		return
	})
	if nil != err {
		panic(err)
	}
}

func AddToTree(name string) {
	namePart := strings.Split(name, "/")
	currNode := LmdbTree
	for ith, k := range namePart {
		if ith == len(namePart)-1 {
			if _, ok := currNode.Children[k]; !ok {
				currNode.Children[k] = &Node{Name: k, IsFile: true, Parent: currNode, Children: make(map[string]*Node)}
			} else {
				currNode.Children[k].IsFile = true
			}
		} else {
			if _, ok := currNode.Children[k]; !ok {
				currNode.Children[k] = &Node{Name: k, IsFile: false, Parent: currNode, Children: make(map[string]*Node)}
			}
			currNode = currNode.Children[k]
		}
	}
}

func GetNeighborFolder(node *Node) (pre, nxt string) {
	parent := node.Parent
	if nil == node.Parent {
		return
	}
	content, _ := GetFolderContent(parent)
	currentName := GetPath(node)
	for i, val := range content.Folders {
		if val == currentName {
			if i-1 >= 0 {
				pre = content.Folders[i-1]
			}
			if i+1 < len(content.Folders) {
				nxt = content.Folders[i+1]
			}
			return
		}
	}
	return
}

func GetFolderContent(root *Node) (content util.FolderContent, err error) {
	content = util.FolderContent{
		Name:    root.Name,
		Folders: []string{},
		Files:   []string{},
	}
	basePath := GetPath(root)
	for _, v := range root.Children {
		pa := path.Join(basePath, v.Name)
		if v.IsFile {
			content.Files = append(content.Files, pa)
		} else {
			content.Folders = append(content.Folders, pa)
		}
	}

	// Sort to keep a static order
	sort.Strings(content.Folders)
	sort.Strings(content.Files)

	return content, nil
}

func GetPath(node *Node) (path string) {
	if node.Parent == nil {
		return "/" + node.Name
	}
	return GetPath(node.Parent) + "/" + node.Name
}

func GetNode(path string) (*Node, error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	currNode := LmdbTree
	if "" != path {
		namePart := strings.Split(path, "/")
		for _, k := range namePart {
			if _, ok := currNode.Children[k]; !ok {
				return nil, nil
			} else {
				currNode = currNode.Children[k]
			}
		}
	}
	return currNode, nil
}
