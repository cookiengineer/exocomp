package utils

import "fmt"
import "strings"

func FormatFileBuffer(raw string) ([]byte, error) {

	formatted := make([]byte, 0)
	err       := error(nil)

	lines := strings.Split(raw, "\n")

	for l := 0; l < len(lines); l++ {

		line := lines[l]

		for p := 0; p < len(line); p++ {

			chr := rune(line[p])

			if line[p] <= 127 {

				switch {
				// digits
				case chr >= '0' && chr <= '9':
					formatted = append(formatted, byte(chr))
				// letters
				case chr >= 'A' && chr <= 'Z':
					formatted = append(formatted, byte(chr))
				case chr >= 'a' && chr <= 'z':
					formatted = append(formatted, byte(chr))
				// whitespaces
				case chr == ' ' || chr == '\n' || chr == '\t':
					formatted = append(formatted, byte(chr))
				// arithmetic operators
				case chr == '+' || chr == '-' || chr == '*' || chr == '/' || chr == '%' || chr == '&':
					formatted = append(formatted, byte(chr))
				// bitwise operators
				case chr == '|' || chr == '^' || chr == '<' || chr == '>':
					formatted = append(formatted, byte(chr))
				// assignment and comparison operators
				case chr == '!' || chr == '=' || chr == ':':
					formatted = append(formatted, byte(chr))
				// delimiters
				case chr == '(' || chr == ')' || chr == '[' || chr == ']' || chr == '{' || chr == '}':
					formatted = append(formatted, byte(chr))
				case chr == ',' || chr == '.' || chr == ';' || chr == ':':
					formatted = append(formatted, byte(chr))
				// string delimiters
				case chr == '"' || chr == '\'' || chr == '\\' || chr == '`':
					formatted = append(formatted, byte(chr))
				// identifier character
				case chr == '_':
					formatted = append(formatted, byte(chr))
				// exocomp gadget call characters
				case chr == '!' || chr == '#':
					formatted = append(formatted, byte(chr))
				// other characters
				case chr == '?' || chr == '~' || chr == '@' || chr == '$':
					formatted = append(formatted, byte(chr))
				default:
					err = fmt.Errorf("Invalid character %s at line %d position %d", string(chr), l+1, p+1)
					break
				}

			} else {
				err = fmt.Errorf("Invalid character %s at line %d position %d", string(chr), l+1, p+1)
				break
			}

		}

		formatted = append(formatted, byte('\n'))

	}

	if err == nil {
		return formatted, nil
	} else {
		return []byte{}, err
	}

}
