package html

import "fmt"
import "strings"

type Element struct {
	Type       string            `json:"type"`
	Text       string            `json:"text"`
	Attributes map[string]string `json:"attributes"`
	Children   []*Element        `json:"children"`
}

func NewElement(typ string) *Element {

	element := &Element{
		Type:       strings.ToLower(strings.TrimSpace(typ)),
		Text:       "",
		Attributes: make(map[string]string),
		Children:   make([]*Element, 0),
	}

	return element

}

func NewTextElement(text string) *Element {

	element := NewElement("#text")
	element.Text = text

	return element

}

func (element *Element) AddChild(child *Element) {

	if element != nil && child != nil {
		element.Children = append(element.Children, child)
	}

}

func (element *Element) GetAttribute(key string) string {

	var result string

	if element != nil {

		value, ok := element.Attributes[key]

		if ok == true {
			result = value
		}

	}

	return result

}

func (element *Element) SetAttribute(key string, value string) {

	if element != nil {
		element.Attributes[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

}

func (element *Element) SetText(value string) {

	if element != nil {
		element.Text = value
	}

}

// textContent returns the raw text of the element and all its children,
// ignoring all markup.
func (element *Element) textContent() string {

	var result string

	if element == nil {
		return result
	}

	if len(element.Children) > 0 {

		for _, child := range element.Children {
			result += child.textContent()
		}

	} else if element.Text != "" {
		result = element.Text
	}

	return result

}

// RenderInline renders the element as a single line of markdown, used inside
// paragraphs, headings, links and list items.
func (element *Element) RenderInline(document *Document) string {

	var result string

	if element == nil {
		return result
	}

	switch element.Type {

	case "#text":
		result = element.Text

	case "a":
		result = fmt.Sprintf("[%s](%s)", strings.TrimSpace(element.renderChildrenInline(document)), resolveURL(document, element.GetAttribute("href")))

	case "b", "strong":
		result = fmt.Sprintf("**%s**", strings.TrimSpace(element.renderChildrenInline(document)))

	case "em", "i":
		result = fmt.Sprintf("*%s*", strings.TrimSpace(element.renderChildrenInline(document)))

	case "code":
		result = fmt.Sprintf("`%s`", element.textContent())

	case "img":
		result = fmt.Sprintf("![%s](%s)", element.GetAttribute("alt"), resolveURL(document, element.GetAttribute("src")))

	case "br":
		result = "\n"

	case "del", "s", "strike":
		result = fmt.Sprintf("~~%s~~", strings.TrimSpace(element.renderChildrenInline(document)))

	default:
		// abbr, span, sub, sup, small, mark, time, u, kbd, var, cite, q
		result = element.renderChildrenInline(document)

	}

	return result

}

func (element *Element) renderChildrenInline(document *Document) string {

	var result string

	if element != nil {

		for _, child := range element.Children {
			result += child.RenderInline(document)
		}

	}

	return result

}

// Render renders the element as a block of markdown.
func (element *Element) Render(document *Document) string {

	var result string

	if element == nil {
		return result
	}

	switch element.Type {

	case "h1", "h2", "h3", "h4", "h5", "h6":

		level := int(element.Type[1] - '0')
		result = strings.Repeat("#", level) + " " + strings.TrimSpace(element.renderChildrenInline(document))

	case "p":
		result = strings.TrimSpace(element.renderChildrenInline(document))

	case "pre":

		text := strings.TrimSpace(element.textContent())

		if text != "" {
			result = "```\n" + text + "\n```"
		}

	case "ul":

		lines := make([]string, 0)

		for _, child := range element.Children {

			if child.Type == "li" {
				lines = append(lines, "- "+strings.TrimSpace(child.renderChildrenInline(document)))
			}

		}

		result = strings.Join(lines, "\n")

	case "ol":

		lines := make([]string, 0)
		index := 1

		for _, child := range element.Children {

			if child.Type == "li" {
				lines = append(lines, fmt.Sprintf("%d. %s", index, strings.TrimSpace(child.renderChildrenInline(document))))
				index++
			}

		}

		result = strings.Join(lines, "\n")

	case "blockquote":

		lines := make([]string, 0)

		for _, child := range element.Children {

			content := child.Render(document)

			for _, line := range strings.Split(content, "\n") {
				lines = append(lines, "> "+line)
			}

		}

		result = strings.Join(lines, "\n")

	case "table":
		result = renderTable(element, document)

	case "hr":
		result = "---"

	case "li":
		result = strings.TrimSpace(element.renderChildrenInline(document))

	default:
		result = element.renderChildrenBlock(document)

	}

	return result

}

func (element *Element) renderChildrenBlock(document *Document) string {

	blocks := make([]string, 0)

	if element != nil {

		for _, child := range element.Children {

			if isBlockElement(child.Type) {

				content := child.Render(document)

				if strings.TrimSpace(content) != "" {
					blocks = append(blocks, content)
				}

			} else {

				inline := child.RenderInline(document)

				if strings.TrimSpace(inline) != "" {
					blocks = append(blocks, inline)
				}

			}

		}

	}

	return strings.Join(blocks, "\n")

}
