package filter

import "github.com/spyzhov/ajson"

type JsonFilter struct {
	jsonPath string
}

func NewJsonFilter(jsonPath string) *JsonFilter {
	if "@DEFAULT" == jsonPath {
		jsonPath = ""
	}
	filter := &JsonFilter{jsonPath: jsonPath}
	return filter
}

func (f *JsonFilter) Filter(fileName string) bool {
	return true
}

func (f *JsonFilter) Extract(line string) []string {
	if "" == f.jsonPath {
		return []string{line}
	}

	var result []string

	root, err := ajson.Unmarshal([]byte(line))
	if nil != err {
		return result
	}
	nodes, err := root.JSONPath(f.jsonPath)
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

	return result
}
