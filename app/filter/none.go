package filter

type NoFilter struct {
}

func NewNoFilter() *NoFilter {
	filter := &NoFilter{}
	return filter
}

func (f *NoFilter) Filter(fileName string) bool {
	return true
}

func (f *NoFilter) Extract(line string) []string {
	return []string{line}
}
