package lmdb_handler

import (
	"path"
	"sort"
	"strings"

	"github.com/bmatsuo/lmdb-go/lmdb"
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
	_, folders, _ := GetFolderContent(parent)
	currentName := GetPath(node)
	for i, val := range folders {
		if val == currentName {
			if i-1 >= 0 {
				pre = folders[i-1]
			}
			if i+1 < len(folders) {
				nxt = folders[i+1]
			}
			return
		}
	}
	return
}

func GetFolderContent(root *Node) (folders []string, files []string, err error) {
	basePath := GetPath(root)
	for _, v := range root.Children {
		pa := path.Join(basePath, v.Name)
		if v.IsFile {
			files = append(files, pa)
		} else {
			folders = append(folders, pa)
		}
	}

	// Sort to keep a static order
	sort.Strings(folders)
	sort.Strings(files)

	return folders, files, nil
}

func GetPath(node *Node) (path string) {
	if node.Parent == nil {
		return "/" + node.Name
	}
	return GetPath(node.Parent) + "/" + node.Name
}
