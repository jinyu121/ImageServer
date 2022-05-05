package util

import (
	"net"
	"path/filepath"
	"sort"
	"strings"

	"haoyu.love/ImageServer/app"
)

func Pagination(pageSize, pageNum int, folders, files []string) ([]string, []string, int, int, int, int) {
	itemsCount := len(folders) + len(files)
	if 0 == itemsCount {
		return folders, files, 1, 1, -1, -1
	}

	pageNumMax := (itemsCount + pageSize - 1) / pageSize

	if pageNum < 1 {
		pageNum = 1
	}
	if pageNum > pageNumMax {
		pageNum = pageNumMax
	}
	pageNumOffsetStart := (pageNum - 1) * (*app.PageSize)
	pageNumOffsetEnd := pageNum * (*app.PageSize)
	if pageNumOffsetEnd > itemsCount {
		pageNumOffsetEnd = itemsCount
	}

	pagePrev := -1
	pageNext := -1
	if pageNumMax > 1 {
		if pageNum > 1 {
			pagePrev = pageNum - 1
		}
		if pageNum < pageNumMax {
			pageNext = pageNum + 1
		}
	}

	if len(folders) > 0 {
		if len(folders) < pageNumOffsetStart {
			pageNumOffsetStart -= len(folders)
			folders = []string{}
		} else {
			tmpEnd := pageNumOffsetEnd
			if len(folders) < pageNumOffsetEnd {
				tmpEnd = len(folders)
			}
			folders = folders[pageNumOffsetStart:tmpEnd]
		}
		pageNumOffsetEnd -= len(folders)
	}
	if len(files) > 0 {
		files = files[pageNumOffsetStart:pageNumOffsetEnd]
	}

	return folders, files, pageNum, pageNumMax, pagePrev, pageNext
}

func AlignStringArrays(data [][]string) [][]string {
	// How many arrays
	n := len(data)
	// Golang doesn't have dataset of set, so we have to use map
	filesSet := make([]map[string]string, n)
	// fileSet stores all the files
	fileSet := make(map[string]bool)
	// Record each array
	for i, dataList := range data {
		filesSet[i] = make(map[string]string)
		for _, f := range dataList {
			name := filepath.Base(f)
			fileSet[name] = true
			filesSet[i][name] = f
		}
	}
	// Now we can get a non-duplicated file list
	var i = 0
	fileList := make([]string, len(fileSet))
	for k := range fileSet {
		fileList[i] = k
		i++
	}
	sort.Strings(fileList)

	// Make final result
	result := make([][]string, len(fileSet))
	for i, k := range fileList {
		line := make([]string, n+1)
		line[0] = k
		for j, fileSetItem := range filesSet {
			v, ok := fileSetItem[k]
			if ok {
				line[j+1] = v
			} else {
				line[j+1] = ""
			}
		}
		result[i] = line
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

func RemoveLeft(str string, data []string) []string {
	for i := range data {
		data[i] = strings.TrimPrefix(data[i], str)
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
