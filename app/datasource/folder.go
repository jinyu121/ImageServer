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
	rootAbsolute, _ := filepath.Abs(root)
	ds := &FolderDataSource{Root: rootAbsolute}
	return ds
}

func (ds *FolderDataSource) GetFile(filePath string) ([]byte, error) {
	currentAbs, _ := AbsolutePath(ds.Root, filePath)
	data, err := os.ReadFile(currentAbs)
	return data, err
}

func (ds *FolderDataSource) GetFolder(current string) (content FolderContent, err error) {
	currentAbs, currentRelative := AbsolutePath(ds.Root, current)

	content = FolderContent{
		Name:    currentRelative,
		Folders: []string{},
		Files:   []string{},
	}

	// Ensure the Root is a folder
	rootInfo, err := os.Stat(currentAbs)
	if nil != err {
		return
	}
	if !rootInfo.Mode().IsDir() {
		return
	}

	// Open the folder
	folder, err := os.Open(currentAbs)
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
			content.Folders = append(content.Folders, path.Join(currentRelative, item.Name()))
		} else {
			content.Files = append(content.Files, path.Join(currentRelative, item.Name()))
		}
	}

	// Sort to keep a static order
	sort.Strings(content.Folders)
	sort.Strings(content.Files)

	return
}

func (ds *FolderDataSource) GetNeighbor(current string) (nav *Navigation) {
	nav = &Navigation{}
	if "/" == current || "" == current {
		return
	}

	_, currentRelative := AbsolutePath(ds.Root, current)
	nav.Current = currentRelative
	currentName := path.Base(current)

	baseAbs, baseRelative := AbsolutePath(ds.Root, path.Dir(currentRelative))
	nav.Parent = baseRelative

	folder, err := os.Open(baseAbs)
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
				nav.Prev = path.Join(baseRelative, folders[i-1])
			}
			if i+1 < len(folders) {
				nav.Next = path.Join(baseRelative, folders[i+1])
			}
			return
		}
	}
	return
}

func (ds *FolderDataSource) Stat(filePath string) *FileStat {
	result := &FileStat{
		Exists: false,
		IsFile: false,
	}

	fullPath, _ := AbsolutePath(ds.Root, filePath)
	if !strings.HasPrefix(fullPath, ds.Root) {
		return result
	}

	if fileInfo, err := os.Stat(fullPath); nil == err {
		result.Exists = true
		result.IsFile = !fileInfo.IsDir()
	}

	return result
}
