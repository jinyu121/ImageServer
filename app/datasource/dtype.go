package datasource

import (
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
