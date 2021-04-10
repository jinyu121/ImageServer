package app

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"html/template"
	"path"
	"strings"
)

func GetTemplate(storage embed.FS, fileList ...string) *template.Template {
	tpl, _ := template.New(path.Base(fileList[0])).Funcs(templateFunction).ParseFS(
		storage,
		append(fileList, "static/template/base.tmpl")...,
	)
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
	"stringToMD5": func(s string) string {
		hash := md5.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	},
}
