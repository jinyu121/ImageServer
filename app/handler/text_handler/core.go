package text_handler

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spyzhov/ajson"
)

var (
	Data = make([]string, 0)
)

func initXsv(f *os.File, ext string, column int, jsonP string, bar *progressbar.ProgressBar) {
	csvReader := csv.NewReader(f)
	if ".tsv" == ext {
		csvReader.Comma = '\t'
	}

	for {
		_ = bar.Add(1)
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(rec) < column {
			continue
		}
		Data = append(Data, ExtractByJsonPath([]byte(rec[column]), jsonP)...)
	}
}

func initText(f *os.File, jsonP string, bar *progressbar.ProgressBar) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_ = bar.Add(1)
		Data = append(Data, ExtractByJsonPath(scanner.Bytes(), jsonP)...)
	}
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
			Data = append(Data, s)
		}
	}
	return result
}
