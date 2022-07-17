package app

import (
	"bytes"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app/datasource"
)

func Paginate(content *[]datasource.FolderContent, size int, current int, url string) Pagination {
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

func AlignContent(contents *[]datasource.FolderContent) datasource.FolderContent {
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
	return datasource.FolderContent{Name: "", Folders: folders, Files: files}
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

func ArrayToSet(data map[string]struct{}, arr []string) {
	for _, v := range arr {
		data[v] = struct{}{}
	}
}

func GetIPAddress() []net.IP {
	result := make([]net.IP, 0)

	ifaces, err := net.Interfaces()
	if nil != err {
		return result
	}

	for _, face := range ifaces {
		if addrs, err := face.Addrs(); nil == err {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				default:
					continue
				}

				if !ip.IsUnspecified() &&
					!ip.IsMulticast() &&
					!ip.IsInterfaceLocalMulticast() &&
					!ip.IsLinkLocalMulticast() &&
					!ip.IsLinkLocalUnicast() {
					result = append(result, ip)
				}
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i], result[j]) < 0
	})

	return result
}

func GetCurrentUrl(c *gin.Context) string {
	p := c.Request.URL.Path
	q := c.Request.URL.Query()
	u, _ := url.Parse(p)
	u.RawQuery = q.Encode()
	return u.String()
}

func IsSubFolder(parent string, sub string) (bool, error) {
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
