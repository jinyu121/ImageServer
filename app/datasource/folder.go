package datasource

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type FolderDataSource struct {
	Root string
}

func NewFolderDataSource(root string) *FolderDataSource {
	ds := &FolderDataSource{Root: root}
	return ds
}

func (ds *FolderDataSource) GetFile(filePath string) ([]byte, error) {
	return nil, nil
}

func (ds *FolderDataSource) GetFolder(current string) (content FolderContent, err error) {
	content = FolderContent{
		Name:    current,
		Folders: []string{},
		Files:   []string{},
	}

	// Ensure the Root is a folder
	rootInfo, err := os.Stat(current)
	if nil != err {
		return
	}
	if !rootInfo.Mode().IsDir() {
		return
	}

	// Open the folder
	folder, err := os.Open(current)
	if nil != err {
		return
	}
	defer func() { _ = folder.Close() }()

	// Iter the folder
	fileInfo, err := folder.Readdir(-1)
	if nil != err {
		return
	}

	// Split result into folders and files
	for _, item := range fileInfo {
		// Filter out hidden files
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}

		if item.IsDir() {
			content.Folders = append(content.Folders, path.Join(current, item.Name()))
		} else {
			content.Files = append(content.Files, path.Join(current, item.Name()))
		}
	}

	// Sort to keep a static order
	sort.Strings(content.Folders)
	sort.Strings(content.Files)

	return
}

func (ds *FolderDataSource) GetNeighbor(current string) (pre string, nxt string) {
	basePath := path.Dir(current)
	currentName := path.Base(current)

	folder, err := os.Open(basePath)
	if nil != err {
		return
	}
	defer func() { _ = folder.Close() }()

	fileInfo, err := folder.Readdir(-1)
	if nil != err {
		return
	}

	// Find all folders
	folders := make([]string, 0)
	for _, item := range fileInfo {
		// Filter out hidden files
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}

		if item.IsDir() {
			folders = append(folders, item.Name())
		}
	}

	// Sort to keep a static order
	sort.Strings(folders)

	for i, val := range folders {
		if val == currentName {
			if i-1 >= 0 {
				pre = path.Join(basePath, folders[i-1])
			}
			if i+1 < len(folders) {
				nxt = path.Join(basePath, folders[i+1])
			}
			return
		}
	}
	return
}

func (ds *FolderDataSource) Stat(filePath string) *FileStat {
	fullPath := filepath.Join(ds.Root, filePath)
	fullPath, _ = filepath.Abs(fullPath)
	result := &FileStat{
		Exists: false,
		IsFile: false,
	}

	if fileInfo, err := os.Stat(filePath); nil == err {
		result.Exists = true
		result.IsFile = !fileInfo.IsDir()
	}

	return result
}
