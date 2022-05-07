package folder_handler

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"haoyu.love/ImageServer/app/util"
)

// GetFolderContent gets folders and files from the given directory
func GetFolderContent(root string) (content util.FolderContent, err error) {
	content = util.FolderContent{
		Name:    root,
		Folders: []string{},
		Files:   []string{},
	}

	// Ensure the root is a folder
	rootInfo, err := os.Stat(root)
	if nil != err {
		return
	}
	if !rootInfo.Mode().IsDir() {
		return
	}

	// Open the folder
	folder, err := os.Open(root)
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
			content.Folders = append(content.Folders, path.Join(root, item.Name()))
		} else {
			content.Files = append(content.Files, path.Join(root, item.Name()))
		}
	}

	// Sort to keep a static order
	sort.Strings(content.Folders)
	sort.Strings(content.Files)

	return
}

func GetNeighborFolder(current string) (pre, nxt string) {
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

func IsSubFolder(parent, sub string) (bool, error) {
	up := ".." + string(os.PathSeparator)

	// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
	rel, err := filepath.Rel(parent, sub)
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(rel, up) && rel != ".." {
		return true, nil
	}
	return false, nil
}
