package yaml

import "errors"
import "fmt"
import "strings"

type Parser struct {
	lines []parser_line
	index int
	root  *Node
}

func NewParser(data []byte) *Parser {

	lines  := parse_lines(string(data))
	parser := &Parser{
		index: 0,
		lines: lines,
		root:  nil,
	}

	node, err := parser.ParseObject(0)

	if err == nil {
		parser.root = node
	}

	return parser

}

func (parser *Parser) ParseArray(expected_indent int) (*Node, error) {

	result := &Node{
        Kind:          ArrayNode,
        ArrayChildren: []*Node{},
	}

	for parser.index < len(parser.lines) {

		line := parser.lines[parser.index]

		if line.Indent < expected_indent {
			break
		} else if strings.HasPrefix(line.Text, "- ") == false {
			break
		}

		item_value := strings.TrimPrefix(line.Text, "- ")
		child_node := &Node{
			Kind:  ScalarNode,
			Value: parser.ParseScalar(item_value),
		}

		result.ArrayChildren = append(result.ArrayChildren, child_node)

		parser.index++

	}

	return result, nil

}

func (parser *Parser) ParseMultilineString(parent_indent int) string {

	lines := []string{}

	for parser.index < len(parser.lines) {

		line := parser.lines[parser.index]

		if line.Indent <= parent_indent {
			break
		}

		lines = append(lines, strings.TrimSpace(line.Text))

		parser.index++

	}

	return strings.Join(lines, "\n")

}

func (parser *Parser) ParseObject(expected_indent int) (*Node, error) {

	result := &Node{
		Kind:           ObjectNode,
		ObjectChildren: map[string]*Node{},
	}

	for parser.index < len(parser.lines) {

		line := parser.lines[parser.index]

		if line.Indent < expected_indent {
			break
		} else if line.Indent > expected_indent {
			return nil, ParseError{
				Line:    line.Number,
				Message: fmt.Sprintf("Unexpected indentation %d instead of %d", expected_indent, line.Indent),
			}
		}

		key, value, has_value := parser.SplitKeyValue(line.Text)

		parser.index++

		if has_value == true {

			if value == "|" {

				multiline_value := parser.ParseMultilineString(line.Indent)

				result.ObjectChildren[key] = &Node{
                    Kind:  ScalarNode,
                    Value: multiline_value,
				}

				continue

			} else {

				scalar_value := parser.ParseScalar(value)

				result.ObjectChildren[key] = &Node{
					Kind:  ScalarNode,
					Value: scalar_value,
				}

				continue

			}

		}

		if parser.index < len(parser.lines) {

			next_line := parser.lines[parser.index]

			if strings.HasPrefix(next_line.Text, "- ") {

				array_node, err := parser.ParseArray(next_line.Indent)

				if err == nil {
					result.ObjectChildren[key] = array_node
				} else {
					return nil, err
				}

			} else {

				child_node, err := parser.ParseObject(next_line.Indent)

				if err == nil {
					result.ObjectChildren[key] = child_node
				} else {
					return nil, err
				}

			}

		} else {

			result.ObjectChildren[key] = &Node{
				Kind: ObjectNode,
			}

			break

		}

	}

	return result, nil

}

func (parser *Parser) ParseScalar(value string) string {

	result := strings.TrimSpace(value)

	if strings.HasPrefix(result, "\"") && strings.HasSuffix(result, "\"") {
		result = strings.Trim(result, "\"")
	} else if strings.HasPrefix(result, "'") && strings.HasSuffix(result, "'") {
		result = strings.Trim(result, "'")
	}

	return result

}

func (parser *Parser) Root() (*Node, error) {

	root := parser.root

	if root != nil {
		return root, nil
	}

	return nil, errors.New("missing document root node")

}

func (parser *Parser) SplitKeyValue(line string) (string, string, bool) {

	colon_index := strings.Index(line, ":")

	if colon_index != -1 {

		key   := strings.TrimSpace(line[:colon_index])
		value := strings.TrimSpace(line[colon_index+1:])

		if value != "" {
			return key, value, true
		} else {
			return key, "", false
		}

	} else {
		return "", "", false
	}

}
