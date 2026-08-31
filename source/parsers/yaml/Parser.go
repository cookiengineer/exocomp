package yaml

import "errors"
import "fmt"
import "strings"

type Parser struct {
	lines []parser_line
	index int
	root  *Node
	err   error
}

func NewParser(data []byte) *Parser {

	lines := parse_lines(string(data))
	parser := &Parser{
		index: 0,
		lines: lines,
		root:  nil,
		err:   nil,
	}

	node, err := parser.ParseObject(0)

	if err != nil {
		parser.err = err
	} else {
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
		} else if line.Indent > expected_indent {
			return nil, ParseError{
				Line:    line.Number,
				Message: fmt.Sprintf("Unexpected indentation %d instead of %d", line.Indent, expected_indent),
			}
		} else if isSequenceItem(line.Text) == false {
			break
		}

		child, err := parser.parseSequenceEntry(expected_indent)

		if err != nil {
			return nil, err
		}

		result.ArrayChildren = append(result.ArrayChildren, child)

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
				Message: fmt.Sprintf("Unexpected indentation %d instead of %d", line.Indent, expected_indent),
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

			}

			result.ObjectChildren[key] = &Node{
				Kind:  ScalarNode,
				Value: parser.ParseScalar(value),
			}

			continue

		}

		child, err := parser.parseNestedBlock(line.Indent)

		if err != nil {
			return nil, err
		}

		if child == nil {
			child = &Node{
				Kind:           ObjectNode,
				ObjectChildren: map[string]*Node{},
			}
		}

		result.ObjectChildren[key] = child

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

	if parser.root != nil {
		return parser.root, nil
	}

	if parser.err != nil {
		return nil, parser.err
	}

	return nil, errors.New("missing document root node")

}

func (parser *Parser) SplitKeyValue(line string) (string, string, bool) {

	line = strings.TrimRight(line, " \t")

	for index := 0; index < len(line); index++ {

		if line[index] != ':' {
			continue
		}

		if index+1 == len(line) {
			return strings.TrimSpace(line[:index]), "", false
		}

		if line[index+1] == ' ' || line[index+1] == '\t' {

			value := strings.TrimSpace(line[index+1:])

			if value != "" {
				return strings.TrimSpace(line[:index]), value, true
			}

			return strings.TrimSpace(line[:index]), "", false

		}

	}

	return "", "", false

}

func (parser *Parser) parseNestedBlock(parent_indent int) (*Node, error) {

	if parser.index >= len(parser.lines) {
		return nil, nil
	}

	next_line := parser.lines[parser.index]

	if next_line.Indent <= parent_indent {
		return nil, nil
	}

	if isSequenceItem(next_line.Text) {
		return parser.ParseArray(next_line.Indent)
	}

	return parser.ParseObject(next_line.Indent)

}

func (parser *Parser) parseSequenceEntry(expected_indent int) (*Node, error) {

	line := parser.lines[parser.index]
	content := strings.TrimSpace(strings.TrimPrefix(line.Text, "-"))

	parser.index++

	return parser.parseSequenceItem(expected_indent, content)

}

func (parser *Parser) parseSequenceItem(expected_indent int, content string) (*Node, error) {

	if content == "" {

		child, err := parser.parseNestedBlock(expected_indent)

		if err != nil {
			return nil, err
		}

		if child == nil {
			return &Node{
				Kind:  ScalarNode,
				Value: "",
			}, nil
		}

		return child, nil

	}

	if isSequenceItem(content) {
		return parser.parseInlineSequence(expected_indent, content)
	}

	key, value, has_value := parser.SplitKeyValue(content)

	if key == "" && has_value == false {
		return &Node{
			Kind:  ScalarNode,
			Value: parser.ParseScalar(content),
		}, nil
	}

	object := &Node{
		Kind:           ObjectNode,
		ObjectChildren: map[string]*Node{},
	}

	if has_value == true {

		object.ObjectChildren[key] = &Node{
			Kind:  ScalarNode,
			Value: parser.ParseScalar(value),
		}

	} else {

		child, err := parser.parseNestedBlock(expected_indent)

		if err != nil {
			return nil, err
		}

		if child == nil {
			child = &Node{
				Kind:           ObjectNode,
				ObjectChildren: map[string]*Node{},
			}
		}

		object.ObjectChildren[key] = child

	}

	if err := parser.consumeMappingFields(expected_indent, object); err != nil {
		return nil, err
	}

	return object, nil

}

func (parser *Parser) parseInlineSequence(expected_indent int, content string) (*Node, error) {

	result := &Node{
		Kind:          ArrayNode,
		ArrayChildren: []*Node{},
	}

	inner_indent := expected_indent + 2
	first_content := strings.TrimSpace(strings.TrimPrefix(content, "-"))

	first_item, err := parser.parseSequenceItem(inner_indent, first_content)

	if err != nil {
		return nil, err
	}

	result.ArrayChildren = append(result.ArrayChildren, first_item)

	for parser.index < len(parser.lines) {

		line := parser.lines[parser.index]

		if line.Indent < inner_indent {
			break
		} else if line.Indent > inner_indent {
			return nil, ParseError{
				Line:    line.Number,
				Message: fmt.Sprintf("Unexpected indentation %d instead of %d", line.Indent, inner_indent),
			}
		} else if isSequenceItem(line.Text) == false {
			break
		}

		item, err := parser.parseSequenceEntry(inner_indent)

		if err != nil {
			return nil, err
		}

		result.ArrayChildren = append(result.ArrayChildren, item)

	}

	return result, nil

}

func (parser *Parser) consumeMappingFields(expected_indent int, object *Node) error {

	for parser.index < len(parser.lines) {

		line := parser.lines[parser.index]

		if line.Indent <= expected_indent {
			break
		}

		if isSequenceItem(line.Text) {
			break
		}

		child, err := parser.ParseObject(line.Indent)

		if err != nil {
			return err
		}

		for key, value := range child.ObjectChildren {
			object.ObjectChildren[key] = value
		}

	}

	return nil

}

func isSequenceItem(text string) bool {

	if strings.HasPrefix(text, "- ") {
		return true
	}

	return text == "-"

}
