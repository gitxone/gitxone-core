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
		v, marked := ValueTypes[suggestion]
		if !marked {
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
		multiple := false
		for _, token := range tokens {
			if token == "" {
				continue
			}
			status, marked := done[serial]
			if marked && status.state == DONE {
				break
			}
			if marked && status.state == REMOVED {
				break
			}
			if marked && status.state == OPTIONAL {
				if token[0] != '-' {
					status.state = DONE
					done[serial] = status
					break
				}
			}
			if marked && status.state == SHORT {
				if token[0] == '-' {
					status.state = INVALID
					done[serial] = status
					break
				}
				if multiple {
					status.state = WAITING
				} else {
					status.state = DONE
				}
				done[serial] = status
			}

			value, multiple, found := option.getValue(token)
			if found {
				for _, s := range option.InvalidSerials {
					done[s] = Status{REMOVED, "", false}
				}
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
		status, marked := done[serial]
		if marked && status.state == SHORT {
			result = append(result, formatStringAsValue(status.value))
			return result
		}
	}
	// removing loop
	group := 0
	trailing := false
	for serial, option := range options {
		if group != option.Group {
			group = option.Group
			trailing = false
		}
		status, marked := done[serial]
		if !trailing && !marked && !option.Optional {
			trailing = true
			continue
		}
		if trailing {
			status.state = REMOVED
			done[serial] = status
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
			result = append(result, formatStringAsValue(v))
		}
	}
	return result
}

func formatStringAsValue(value string) string {
	if value == "" {
		return ""
	}
	if value[0] == '-' || value[0] == '=' {
		return value[1:]
	}
	values := make([]string, 0)
	for _, s := range strings.Split(value, " ") {
		values = append(values, fmt.Sprintf("<%s>", s))
	}
	return strings.Join(values, " ")
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
