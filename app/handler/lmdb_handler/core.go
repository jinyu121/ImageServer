package lmdb_handler

import (
	"github.com/bmatsuo/lmdb-go/lmdb"
	"strings"
)

type Node struct {
	Name     string
	Leaf     bool
	Children map[string]*Node
}

var (
	LmdbEnv  *lmdb.Env
	LmdbDBI  lmdb.DBI
	LmdbTree = &Node{Name: "/", Leaf: false, Children: make(map[string]*Node)}
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
				currNode.Children[k] = &Node{Name: k, Leaf: true, Children: make(map[string]*Node)}
			} else {
				currNode.Children[k].Leaf = true
			}
		} else {
			if _, ok := currNode.Children[k]; !ok {
				currNode.Children[k] = &Node{Name: k, Leaf: false, Children: make(map[string]*Node)}
			}
			currNode = currNode.Children[k]
		}
	}
}
