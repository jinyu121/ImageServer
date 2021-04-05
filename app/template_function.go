package app

import (
	"html/template"
	"path"
	"strings"
)

var TemplateFunction = template.FuncMap{
	"pathToName": func(p string) string {
		return path.Base(p)
	},
	"lastOne": func(arr []interface{}) interface{} {
		if 0 == len(arr) {
			return nil
		}
		return arr[len(arr)-1]
	},
	"breadCrumb": func(root string) []string {
		crumb := make([]string, 0)
		root = strings.Trim(root, "/")
		if "" != root {
			sps := strings.Split(root, "/")
			for i := range sps {
				crumb = append(crumb, "/"+strings.Join(sps[:i+1], "/"))
			}
		}
		return crumb
	},
}
