package app

import (
	"net/url"
	"strconv"

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
