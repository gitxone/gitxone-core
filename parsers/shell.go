package parsers

const doubleQuot = '"'
const singleQuot = '\''
const backSlash = '\\'
const space = ' '

// ParseCommand splits a string to tokens for exec.Command
func ParseCommand(s string) []string {
	var prevChar rune
	var quot rune
	var token []rune
	var tokens []string

	for _, char := range s {
		switch char {
		case doubleQuot:
			switch quot {
			case 0:
				quot = doubleQuot
			case doubleQuot:
				if prevChar == backSlash {
					token = token[:len(token)-1]
					token = append(token, char)
				} else {
					quot = 0
				}
			default:
				token = append(token, char)
			}
		case singleQuot:
			switch quot {
			case 0:
				quot = singleQuot
			case singleQuot:
				if prevChar == backSlash {
					token = token[:len(token)-1]
					token = append(token, char)
				} else {
					quot = 0
				}
			default:
				token = append(token, char)
			}
		case space:
			if quot == 0 {
				if prevChar != space {
					tokens = append(tokens, string(token))
					token = []rune("")
				}
			} else {
				token = append(token, char)
			}
		default:
			token = append(token, char)
		}
		prevChar = char
	}

	tokens = append(tokens, string(token))
	return tokens
}
