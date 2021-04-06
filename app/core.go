package app

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
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
	defer folder.Close()

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

// GetTextContent gets non-empty lines from text file
func GetTextContent(root string) (lines []string, err error) {
	// Ensure the root is a folder
	rootInfo, err := os.Stat(root)
	if nil != err {
		return
	}
	if rootInfo.Mode().IsDir() {
		return
	}

	// Open the file
	f, err := os.Open(root)
	if err != nil {
		return
	}
	defer f.Close()

	// Read content
	scanner := bufio.NewScanner(f)
	var text string
	for scanner.Scan() {
		text = strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			lines = append(lines, text)
		}
	}

	return
}

func GetNeighborFolder(current, root string, offset int) (pre, nxt string) {
	if strings.TrimRight(current, "/") == strings.TrimRight(root, "/") {
		return
	}
	baseFolder := path.Dir(current)
	currentFolder := path.Base(current)

	folder, err := os.Open(baseFolder)
	if nil != err {
		return
	}
	defer folder.Close()

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
		if val == currentFolder {
			if i-offset >= 0 {
				pre = strings.TrimPrefix(path.Join(baseFolder, folders[i-offset]), root)
			}
			if i+offset < len(folders) {
				nxt = strings.TrimPrefix(path.Join(baseFolder, folders[i+offset]), root)
			}
			return
		}
	}
	return
}

func IsImageFile(file string) bool {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".jpg", ".jpeg", ".png", ".bmp", ".gif":
		return true
	}
	return false
}

func IsVideoFile(file string) bool {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".mp4", ".avi":
		return true
	}
	return false
}

func FilterFiles(files []string, fn func(string) bool) []string {
	result := make([]string, 0)
	for _, val := range files {
		if fn(val) {
			result = append(result, val)
		}
	}
	return result
}

func RemoveLeft(data []string, str string) []string {
	for i := range data {
		data[i] = strings.TrimPrefix(data[i], str)
	}
	return data
}
