package core

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
		if strings.HasPrefix(filepath.Base(item.Name()), ".") {
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

func isImageFile(file string) bool {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".jpg", ".jpeg", ".png", ".bmp", ".gif":
		return true
	}
	return false
}

func isVideoFile(file string) bool {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".mp4", ".avi":
		return true
	}
	return false
}

func FilterImageFiles(files []string) []string {
	n := 0
	for _, val := range files {
		if isImageFile(val) || isVideoFile(val) {
			files[n] = val
			n++
		}
	}
	return files[:n]
}

func RemoveLeft(data []string, str string) []string {
	for i := range data {
		data[i] = strings.TrimLeft(data[i], str)
	}
	return data
}
