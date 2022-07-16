package util

import (
	"net/url"
	"strconv"
	"strings"
)

type FolderContent struct {
	Name    string
	Folders []string
	Files   []string
}

func (f *FolderContent) FilterTargetFile(target map[string]struct{}) {
	if len(target) > 0 {
		result := make([]string, 0)
		for _, file := range f.Files {
			if IsTargetFileM(file, target) {
				result = append(result, file)
			}
		}
		f.Files = result
	}
}

func (f *FolderContent) RemovePrefix(str string) {
	f.Name = strings.TrimPrefix(f.Name, str)
	if "" == f.Name {
		f.Name = "/"
	}
	f.Folders = RemoveLeft(str, f.Folders, false)
	f.Files = RemoveLeft(str, f.Files, false)
}

type Pagination struct {
	Current int
	Prev    int
	Next    int
	Total   int
	Size    int
	Url     string
	Content *[]FolderContent
}

func (p Pagination) URLPrev() string {
	if p.Prev <= 0 {
		return "#"
	}
	return p.offset(p.Prev)
}

func (p Pagination) URLNext() string {
	if p.Next <= 0 {
		return "#"
	}
	return p.offset(p.Next)
}

func (p Pagination) offset(page int) string {
	u, _ := url.Parse(p.Url)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Set("p", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}

type Navigation struct {
	Current string
	Prev    string
	Next    string
	Parent  string
}

func (f *Navigation) RemovePrefix(str string) {
	f.Current = strings.TrimPrefix(f.Current, str)
	f.Prev = strings.TrimPrefix(f.Prev, str)
	f.Next = strings.TrimPrefix(f.Next, str)
	f.Parent = strings.TrimPrefix(f.Parent, str)
	if "" == f.Parent && "" != f.Current {
		f.Parent = "/"
	}
}
