package app

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var TemplateFunction = template.FuncMap{
	"pathToName": func(p string) string {
		return path.Base(p)
	},
	"lastOne": func(arr []interface{}) interface{} {
		if len(arr) == 0 {
			return nil
		}
		return arr[len(arr)-1]
	},
	"breadCrumb": func(root string) []string {
		crumb := make([]string, 0)
		root = strings.Trim(root, "/")
		if root != "" {
			sps := strings.Split(root, "/")
			for i := range sps {
				crumb = append(crumb, "/"+strings.Join(sps[:i+1], "/"))
			}
		}
		return crumb
	},
	"stringToMD5": func(s string) string {
		hash := md5.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	},
	"pageNum": func(rawUrl string, page int) string {
		if page <= 0 {
			return "#"
		}

		u, _ := url.Parse(rawUrl)
		q, _ := url.ParseQuery(u.RawQuery)
		q.Set("p", strconv.Itoa(page))
		u.RawQuery = q.Encode()
		return u.String()
	},
}
