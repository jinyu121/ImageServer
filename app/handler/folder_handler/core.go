package folder_handler

import (
	"os"
	"path"
	"sort"
	"strings"

	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
)

// GetFolderContent gets folders and files from the given directory
func GetFolderContent(root string) (folders []string, files []string, err error) {
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
			folders = append(folders, path.Join(root, item.Name()))
		} else {
			files = append(files, path.Join(root, item.Name()))
		}
	}

	// Sort to keep a static order
	sort.Strings(folders)
	sort.Strings(files)

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

func FilterTargetFile(files []string) (result []string) {
	result = make([]string, 0)
	for _, file := range files {
		if util.IsTargetFileM(file, app.FileExtension) {
			result = append(result, file)
		}
	}
	return result
}
