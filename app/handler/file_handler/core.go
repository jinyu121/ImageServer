package file_handler

import (
	"bufio"
	"os"
	"strings"
)

// GetTextContent gets non-empty lines from text file
func GetTextContent(root string) (lines []string, err error) {
	// Ensure the root is a folder
	rootInfo, err := os.Stat(root)
	if nil != err {
		return
	}
	if rootInfo.Mode().IsDir() {
		return
	}

	// Open the file
	f, err := os.Open(root)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	// Read content
	scanner := bufio.NewScanner(f)
	var text string
	for scanner.Scan() {
		text = strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			lines = append(lines, text)
		}
	}

	return
}
