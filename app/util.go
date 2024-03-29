package app

import (
	"bytes"
	"net"
	"net/url"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app/datasource"
)

// Paginate paginates the given content in place, and returns the Pagination instance
func Paginate(
	ref *datasource.FolderContent,
	contents *[]datasource.FolderContent,
	size int, current int, url string) Pagination {

	page := Pagination{Current: 1, Prev: -1, Next: -1, Total: 1, Size: size, Toc: ref, Url: url}

	numFolders, numFiles := len(ref.Folders), len(ref.Files)

	if size <= 0 {
		return page
	}
	// Ensure all element are not empty
	itemsCount := numFolders + numFiles
	if itemsCount == 0 {
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
			ref.Folders = ref.Folders[tmpStart:tmpEnd]
		} else {
			tmpStart := offsetStart
			tmpEnd := numFolders
			ref.Folders = ref.Folders[tmpStart:tmpEnd]

			tmpStart = 0
			tmpEnd = offsetEnd - numFolders
			ref.Files = ref.Files[tmpStart:tmpEnd]
		}
	} else {
		tmpStart := offsetStart - numFolders
		tmpEnd := offsetEnd - numFolders
		ref.Folders = make([]string, 0)
		ref.Files = ref.Files[tmpStart:tmpEnd]
	}

	// Align content
	AlignContent(contents, ref)
	page.Content = contents

	return page
}

// DeduplicateFolderContent merges the content of multiple FolderContent instances
func DeduplicateFolderContent(contents *[]datasource.FolderContent) datasource.FolderContent {
	// Deduplicate
	folderSet := make(map[string]struct{})
	fileSet := make(map[string]struct{})
	for _, content := range *contents {
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

	return datasource.FolderContent{Name: "", Folders: folders, Files: files}
}

// AlignContent aligns the content of folders and files in-place according to the given content
func AlignContent(contents *[]datasource.FolderContent, ref *datasource.FolderContent) {
	contents_ := *contents

	// Align
	for i := range contents_ {
		contents_[i].Folders = align(contents_[i].Folders, ref.Folders)
		contents_[i].Files = align(contents_[i].Files, ref.Files)
	}

}

// align the array to the given array of strings.
// If something is not in the given array, a blank will be added in that place.
func align(items, ref []string) []string {
	result := make([]string, len(ref))
	tmp := make(map[string]string)
	for _, item := range items {
		name := filepath.Base(item)
		tmp[name] = item
	}
	for i, item := range ref {
		name := filepath.Base(item)
		if val, ok := tmp[name]; ok {
			result[i] = val
		} else {
			result[i] = ""
		}
	}
	return result
}

func GetIPAddress() []net.IP {
	result := make([]net.IP, 0)

	iFaces, err := net.Interfaces()
	if nil != err {
		return result
	}

	for _, face := range iFaces {
		if addresses, err := face.Addrs(); nil == err {
			for _, addr := range addresses {
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
