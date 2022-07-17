package datasource

import (
	"path/filepath"
	"strings"
)

type FileStat struct {
	Exists bool
	IsFile bool
}

type Node struct {
	Name     string
	IsFile   bool
	Parent   *Node
	Children map[string]*Node
}

func NewRoot() *Node {
	node := &Node{Name: "", IsFile: false, Parent: nil, Children: make(map[string]*Node)}
	return node
}

func (node *Node) Add(name string) {
	namePart := strings.Split(name, "/")
	currNode := node
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

func (node *Node) GetChild(name string) (*Node, error) {
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	current := node
	if "" != name {
		namePart := strings.Split(name, "/")
		for _, k := range namePart {
			if _, ok := current.Children[k]; !ok {
				return nil, nil
			} else {
				current = current.Children[k]
			}
		}
	}
	return current, nil
}

func (node *Node) GetAbsolutePath() (path string) {
	if node.Parent == nil {
		return node.Name
	}
	return node.Parent.GetAbsolutePath() + "/" + node.Name
}

func IsTargetFileM(file string, target ...map[string]struct{}) bool {
	ext := strings.ToLower(filepath.Ext(file))
	for _, t := range target {
		if _, ok := t[ext]; ok {
			return true
		}
	}
	return false
}

func RemoveLeft(str string, data []string, nonEmpty bool) []string {
	for i := range data {
		data[i] = strings.TrimPrefix(data[i], str)
		if "" == data[i] && nonEmpty {
			data[i] = "/"
		}
	}
	return data
}

func AbsolutePath(prefix, relative string) (string, string) {
	relativePart := strings.Trim(strings.TrimSpace(relative), "/")

	absolute := filepath.Join(prefix, relativePart)
	absolute, _ = filepath.Abs(absolute)
	relativeNew := "/" + relativePart

	return absolute, relativeNew
}
