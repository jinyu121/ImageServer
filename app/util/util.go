package util

import (
	"net"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func Paginate(content *[]FolderContent, size int, current int, url string) Pagination {
	page := Pagination{Current: 1, Prev: -1, Next: -1, Total: 1, Size: size, Content: content, Url: url}
	content_ := *content
	// Totally empty
	if 0 == len(content_) {
		return page
	}
	// Ensure all element are all equal in size
	numFolders, numFiles := len(content_[0].Folders), len(content_[0].Files)
	for _, item := range content_ {
		if len(item.Folders) != numFolders || len(item.Files) != numFiles {
			return page
		}
	}

	if size <= 0 {
		return page
	}
	// Ensure all element are not empty
	itemsCount := numFolders + numFiles
	if 0 == itemsCount {
		return page
	}

	// Calculate pagination
	page.Current = current
	page.Total = (itemsCount + size - 1) / size

	if page.Current < 1 {
		page.Current = 1
	}
	if page.Current > page.Total {
		page.Current = page.Total
	}
	offsetStart := (page.Current - 1) * page.Size
	offsetEnd := page.Current * page.Size
	if offsetEnd > itemsCount {
		offsetEnd = itemsCount
	}

	if page.Total > 1 {
		if page.Current > 1 {
			page.Prev = page.Current - 1
		}
		if page.Current < page.Total {
			page.Next = page.Current + 1
		}
	}

	// Limit folders and files
	if offsetStart < numFolders {
		if offsetEnd < numFolders {
			tmpStart := offsetStart
			tmpEnd := offsetEnd
			for i := range content_ {
				content_[i].Folders = content_[i].Folders[tmpStart:tmpEnd]
			}
		} else {
			tmpStart := offsetStart
			tmpEnd := numFolders
			for i := range content_ {
				content_[i].Folders = content_[i].Folders[tmpStart:tmpEnd]
			}
			tmpStart = 0
			tmpEnd = offsetEnd - numFolders
			for i := range content_ {
				content_[i].Files = content_[i].Files[tmpStart:tmpEnd]
			}
		}
	} else {
		tmpStart := offsetStart
		tmpEnd := offsetEnd
		for i := range content_ {
			content_[i].Files = content_[i].Files[tmpStart:tmpEnd]
		}
	}

	return page
}

func AlignContent(contents *[]FolderContent) FolderContent {
	contents_ := *contents
	n := len(contents_)
	if n <= 1 {
		return contents_[0]
	}

	// Deduplicate
	folderSet := make(map[string]struct{})
	fileSet := make(map[string]struct{})
	for _, content := range contents_ {
		for _, folder := range content.Folders {
			name := filepath.Base(folder)
			folderSet[name] = struct{}{}
		}
		for _, file := range content.Files {
			name := filepath.Base(file)
			fileSet[name] = struct{}{}
		}
	}

	// Sort
	folders := make([]string, 0, len(folderSet))
	for folder := range folderSet {
		folders = append(folders, folder)
	}
	sort.Strings(folders)

	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)

	// Align
	for i := range contents_ {
		contents_[i].Folders = align(contents_[i].Folders, folders)
		contents_[i].Files = align(contents_[i].Files, files)
	}
	return FolderContent{Name: "", Folders: folders, Files: files}
}

func align(items, total []string) []string {
	result := make([]string, len(total))
	tmp := make(map[string]string)
	for _, item := range items {
		name := filepath.Base(item)
		tmp[name] = item
	}
	for i, item := range total {
		name := filepath.Base(item)
		if val, ok := tmp[name]; ok {
			result[i] = val
		} else {
			result[i] = ""
		}
	}
	return result
}

func FilterItems(items []string, fn func(string) bool) []string {
	result := make([]string, 0)
	for _, val := range items {
		if fn(val) {
			result = append(result, val)
		}
	}
	return result
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

func IsTargetFileL(file string, target ...[]string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	for _, t := range target {
		for _, v := range t {
			if ext == v {
				return true
			}
		}
	}
	return false
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

func ArrayToSet(data map[string]struct{}, arr []string) {
	for _, v := range arr {
		data[v] = struct{}{}
	}
}

func GetIPAddress() []string {
	result := make([]string, 0)

	ifaces, err := net.Interfaces()
	if nil != err {
		return result
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if nil != err {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() {
				result = append(result, ip.String())
			}
		}
	}

	return result
}

func GetCurrentUrl(c *gin.Context) string {
	p := c.Request.URL.Path
	q := c.Request.URL.Query()
	u, _ := url.Parse(p)
	u.RawQuery = q.Encode()
	return u.String()
}
