package app

import (
	"embed"
	"html/template"
	"path"
	"strings"
)

func GetTemplate(storage embed.FS, fileList ...string) *template.Template {
	var tpl *template.Template
	fileList = append(fileList, "static/template/base.tmpl")
	tpl, _ = template.New(path.Base(fileList[0])).Funcs(templateFunction).ParseFS(storage, fileList...)
	return tpl
}

var templateFunction = template.FuncMap{
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
