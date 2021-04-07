package version

import "fmt"

const (
	CodeName string = "ImageServer"
)

var (
	Version     string = "<undefined>"
	Build       string = "<undefined>"
	VersionLong string = fmt.Sprintf("%s %s (Build %s)", CodeName, Version, Build)
)
