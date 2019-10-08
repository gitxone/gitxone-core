package completion

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Status struct {
	state    int
	value    string
	optional bool
}

const (
	_ int = iota
	DONE
	SHORT
	WAITING
	OPTIONAL
	INVALID
	REMOVED
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

	options, ok := Options[tokens[0]]
	if !ok {
		return []string{}
	}
	result := make([]string, 0)
	for _, suggestion := range suggest(options, tokens[1:]) {
		if suggestion == "" {
			continue
		}
		v, existing := ValueTypes[suggestion]
		if !existing {
			result = append(result, suggestion)
			continue
		}
		switch v {
		case "<file>":
			result = append(result, getFiles(path, lastToken)...)
		default:
			result = append(result, v)
		}
	}
	return result
}

func (t *Option) getValue(key string) (string, bool, bool) {
	index := -1
	for i, k := range t.Keys {
		if k == key {
			index = i
			break
		}
	}
	if index == -1 {
		return "", false, false
	}
	return t.Values[index], t.Multiple[index], true
}

func suggest(options []Option, tokens []string) []string {
	done := make(map[int]Status, 0)

	// marking loop
	for serial, option := range options { //options
		//value := ""
		multiple := false
		for _, token := range tokens {
			if token == "" {
				continue
			}
			status, checked := done[serial]
			if checked && status.state == DONE {
				break
			}
			if checked && status.state == OPTIONAL {
				if token[0] != '-' {
					status.state = DONE
					break
				}
			}
			if checked && status.state == SHORT {
				if token[0] == '-' {
					status.state = INVALID
					break
				}
				if multiple {
					status.state = WAITING
				} else {
					status.state = DONE
				}
			}

			value, multiple, found := option.getValue(token)
			if found {
				if value == "" {
					if multiple {
						done[serial] = Status{WAITING, "", false}
					} else {
						done[serial] = Status{DONE, "", false}
					}
				} else if value[len(value)-1] == '?' {
					done[serial] = Status{OPTIONAL, value[:len(value)-1], true}
				} else {
					done[serial] = Status{SHORT, value, false}
				}
			}
		}
	}

	result := make([]string, 0)
	// extracting short loop
	for serial := range options {
		status, existing := done[serial]
		if existing && status.state == SHORT {
			result = append(result, fmt.Sprintf("<%s>", status.value))
			return result
		}
	}

	// extracting loop
	for serial, option := range options {
		status, _ := done[serial]
		switch status.state {
		case DONE:
			continue
		case INVALID:
			continue
		case REMOVED:
			continue
		case SHORT:
			result = append(result, status.value)
			break
		}
		result = append(result, option.Keys...)
		for index, v := range option.Values {
			if v == "" || option.Keys[index] != "" {
				continue
			}
			if v[len(v)-1] == '?' {
				v = v[:len(v)-1]
			}
			result = append(result, fmt.Sprintf("<%s>", v))
		}
	}
	return result
}

func getFiles(path string, lastToken string) []string {
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
