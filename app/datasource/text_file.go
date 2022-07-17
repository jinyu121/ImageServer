package datasource

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spyzhov/ajson"
)

type TextFileDataSource struct {
	Root   string
	filter string
	column int
	data   []string
}

func NewTextFileDataSource(root string, filter string, column int) *TextFileDataSource {
	ds := &TextFileDataSource{Root: root, filter: filter, column: column}
	ds.read()
	return ds
}

func (ds *TextFileDataSource) GetFile(filePath string) ([]byte, error) {
	return nil, nil
}

func (ds *TextFileDataSource) GetFolder(current string) (content FolderContent, err error) {
	content = FolderContent{
		Name:    "",
		Folders: []string{},
		Files:   ds.data,
	}
	return content, nil
}

func (ds *TextFileDataSource) GetNeighbor(current string) (nav *Navigation) {
	nav = &Navigation{}
	return nav
}

func (ds *TextFileDataSource) read() {
	log.Printf("Scan column %d of file %s", ds.column, ds.Root)
	bar := progressbar.Default(-1, "Scanning")
	callback := func(path string) {
		_ = bar.Add(1)
	}

	ext := strings.ToLower(filepath.Ext(ds.Root))

	f, err := os.Open(ds.Root)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	if ".csv" == ext || ".tsv" == ext {
		ReadXsv(f, ext, ds.column, ds.filter, callback)
	} else {
		ReadText(f, ds.filter, callback)
	}

	log.Printf("Done! %d records read", len(ds.data))
}

func (ds *TextFileDataSource) Stat(filePath string) *FileStat {
	result := &FileStat{
		Exists: false,
		IsFile: false,
	}

	return result
}

func ReadXsv(f *os.File, ext string, column int, jsonP string, callback func(path string)) []string {
	result := make([]string, 0)

	csvReader := csv.NewReader(f)
	if ".tsv" == ext {
		csvReader.Comma = '\t'
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(rec) < column {
			continue
		}
		result = append(result, ExtractByJsonPath([]byte(rec[column]), jsonP)...)

		if nil != callback {
			callback("")
		}
	}

	return result
}

func ReadText(f *os.File, jsonP string, callback func(path string)) []string {
	result := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, ExtractByJsonPath(scanner.Bytes(), jsonP)...)
		if nil != callback {
			callback("")
		}
	}

	return result
}

func ExtractByJsonPath(data []byte, pa string) []string {
	var result []string
	if "" == pa {
		result = append(result, string(data))
		return result
	}
	if "" != pa {
		root, err := ajson.Unmarshal(data)
		if nil != err {
			return result
		}
		nodes, err := root.JSONPath(pa)
		if nil != err {
			return result
		}
		for _, node := range nodes {
			s, err := node.GetString()
			if nil != err {
				continue
			}
			result = append(result, s)
		}
	}
	return result
}
