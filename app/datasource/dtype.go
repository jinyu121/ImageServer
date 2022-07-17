package datasource

type DataSource interface {
	GetFile(path string) ([]byte, error)
	GetFolder(path string) (FolderContent, error)
	GetNeighbor(current string) *Navigation
	Stat(filePath string) *FileStat
}

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

type Navigation struct {
	Current string
	Prev    string
	Next    string
	Parent  string
}
