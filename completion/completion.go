package completion

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Complete(path string, tokens []string) []string {
	lastToken := tokens[len(tokens)-1]
	if len(tokens) < 2 {
		filtered := make([]string, 0)
		for _, c := range gitCommands {
			index := len(lastToken)
			if index <= len(c) && c[0:index] == lastToken {
				filtered = append(filtered, c)
			}
		}
		return filtered
	}

	dirName, bits := filepath.Split(lastToken)
	joinedPath := filepath.Join(path, dirName)

	files, _ := ioutil.ReadDir(joinedPath)
	paths := make([]string, 0)

	for _, f := range files {
		fileName := f.Name()
		p := filepath.Join(dirName, fileName)
		if f.IsDir() {
			p += "/"
		}
		if strings.HasPrefix(fileName, bits) {
			paths = append(paths, p)
		}
	}
	return paths
}
