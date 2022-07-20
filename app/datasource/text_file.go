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
	"haoyu.love/ImageServer/app/filter"
)

type TextFileDataSource struct {
	root string

	filter *filter.Filter
	column int
	data   []string
}

func NewTextFileDataSource(root string, flt *filter.Filter, column int) *TextFileDataSource {
	ds := &TextFileDataSource{root: root, filter: flt, column: column}
	ds.scan()
	return ds
}

func (ds *TextFileDataSource) GetFile(_ string) ([]byte, error) {
	return nil, nil
}

func (ds *TextFileDataSource) GetFolder(_ string) (content FolderContent, err error) {
	content = FolderContent{
		Name:    "",
		Folders: []string{},
		Files:   ds.data,
	}
	return content, nil
}

func (ds *TextFileDataSource) GetNeighbor(_ string) (nav *Navigation) {
	nav = &Navigation{}
	return nav
}

func (ds *TextFileDataSource) scan() {
	log.Printf("Scan column %d of file %s", ds.column, ds.root)
	bar := progressbar.Default(-1, "Scanning")
	callback := func(path string) {
		_ = bar.Add(1)
	}

	ext := strings.ToLower(filepath.Ext(ds.root))

	f, err := os.Open(ds.root)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	if ".csv" == ext || ".tsv" == ext {
		ds.data = ReadXsv(f, ext, ds.column, ds.filter, callback)
	} else {
		ds.data = ReadText(f, ds.filter, callback)
	}

	log.Printf("Done! %d records scan", len(ds.data))
}

func (ds *TextFileDataSource) Stat(filePath string) *FileStat {
	result := &FileStat{
		Exists: false,
		IsFile: false,
	}
	if "" == filePath || "/" == filePath {
		result.Exists = true
	}

	return result
}

func ReadXsv(f *os.File, ext string, column int, flt *filter.Filter, callback func(path string)) []string {
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
		result = append(result, (*flt).Extract(rec[column])...)

		if nil != callback {
			callback("")
		}
	}

	return result
}

func ReadText(f *os.File, flt *filter.Filter, callback func(path string)) []string {
	result := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, (*flt).Extract(scanner.Text())...)
		if nil != callback {
			callback("")
		}
	}

	return result
}
