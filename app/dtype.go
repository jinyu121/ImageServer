package app

import (
	"net/url"
	"strconv"
	"strings"

	"haoyu.love/ImageServer/app/datasource"
)

type Pagination struct {
	Current int
	Prev    int
	Next    int
	Total   int
	Size    int
	Url     string
	Content *[]datasource.FolderContent
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
