package markdown

import "fmt"
import "sort"
import "strings"

var markdownLineWidth = 90

type Element struct {
	Line             int               `json:"line"`
	Type             string            `json:"type"`
	Text             string            `json:"text"`
	Attributes       map[string]string `json:"attributes"`
	Children         []*Element        `json:"children"`
	is_block_element bool              `json:"-"`
}

func NewElement(typ string) *Element {

	var element Element

	element.Line = int(0)
	element.Type = strings.TrimSpace(typ)
	element.Attributes = make(map[string]string)
	element.Children = make([]*Element, 0)
	element.is_block_element = false

	// There's no clean place to do this, as it's part of the CSS specification
	switch typ {

	case "div", "figure", "figcaption", "p":
		element.is_block_element = true
	case "pre":
		element.is_block_element = true
	case "article", "dialog", "footer", "header", "main":
		element.is_block_element = true
	case "blockquote", "hr":
		element.is_block_element = true
	default:
		element.is_block_element = false

	}

	return &element

}

func (element *Element) AddChild(child *Element) {

	element.Children = append(element.Children, child)

	if element.Type == "h1" || element.Type == "h2" || element.Type == "h3" || element.Type == "h4" {
		element.Attributes["id"] = generateId(element)
	}

}

func (element *Element) AddChildren(children []*Element) {

	filtered := make([]*Element, 0)

	for _, child := range children {

		if child != nil {
			filtered = append(filtered, child)
		}

	}

	for _, child := range filtered {
		element.Children = append(element.Children, child)
	}

	if element.Type == "h1" || element.Type == "h2" || element.Type == "h3" || element.Type == "h4" {
		element.Attributes["id"] = generateId(element)
	}

}

func (element *Element) AddText(value string) {

	if element.Text != "" {
		element.Text = element.Text + "\n" + value
	} else {
		element.Text = value
	}

	if element.Type == "h1" || element.Type == "h2" || element.Type == "h3" || element.Type == "h4" {
		element.Attributes["id"] = generateId(element)
	}

}

func (element *Element) getAttributes() []string {

	result := make([]string, 0)

	for name, _ := range element.Attributes {
		result = append(result, name)
	}

	sort.Strings(result)

	return result

}

func (element *Element) GetAttribute(key string) string {

	var result string

	value, ok := element.Attributes[key]

	if ok == true {
		result = value
	}

	return result

}

func (element *Element) getLastChild() *Element {

	var child *Element = nil

	if len(element.Children) > 0 {
		child = element.Children[len(element.Children)-1]
	}

	return child

}

func (element *Element) HasAttribute(key string) bool {

	var result bool

	_, ok := element.Attributes[key]

	if ok == true {
		result = true
	}

	return result

}

func (element *Element) IsBlockElement() bool {
	return element.is_block_element
}

func (element *Element) SetAttribute(key string, value string) {

	if element.Type == "a" && key == "href" {

		if strings.HasPrefix(value, "https://github.com") {
			element.Attributes["class"] = "icon-github"
		} else if strings.HasPrefix(value, "https://gitlab.com") {
			element.Attributes["class"] = "icon-gitlab"
		} else if strings.HasPrefix(value, "https://instagram.com") {
			element.Attributes["class"] = "icon-instagram"
		} else if strings.HasPrefix(value, "https://linkedin.com") {
			element.Attributes["class"] = "icon-linkedin"
		} else if strings.HasPrefix(value, "https://medium.com") {
			element.Attributes["class"] = "icon-medium"
		} else if strings.HasPrefix(value, "https://reddit.com") {
			element.Attributes["class"] = "icon-reddit"
		} else if strings.HasPrefix(value, "https://old.reddit.com") {
			element.Attributes["class"] = "icon-reddit"
		} else if strings.HasPrefix(value, "https://cookie.engineer") {
			// Do Nothing
		} else if strings.HasPrefix(value, "../") {
			// Do Nothing
		} else if strings.HasPrefix(value, "./") {
			// Do Nothing
		} else if strings.HasPrefix(value, "/") {
			// Do Nothing
		} else if strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://") {
			element.Attributes["class"] = "icon-website"
		}

		if strings.HasSuffix(value, ".go") ||
			strings.HasSuffix(value, ".js") ||
			strings.HasSuffix(value, ".mjs") ||
			strings.HasSuffix(value, ".pdf") ||
			strings.HasSuffix(value, ".tar.gz") ||
			strings.HasSuffix(value, ".tar.xz") ||
			strings.HasSuffix(value, ".zip") {
			element.Attributes["class"] = "icon-download"
		} else if strings.HasPrefix(value, "#") {
			element.Attributes["class"] = "icon-section"
		}

		element.Attributes[key] = strings.TrimSpace(value)

	} else if element.Type == "img" && key == "src" {

		element.Attributes[key] = strings.TrimSpace(value)

	} else {

		element.Attributes[key] = strings.TrimSpace(value)

	}

}

func (element *Element) SetChildren(children []*Element) {

	filtered := make([]*Element, 0)

	for _, child := range children {

		if child != nil {
			filtered = append(filtered, child)
		}

	}

	element.Children = filtered

	if element.Type == "h1" || element.Type == "h2" || element.Type == "h3" || element.Type == "h4" {
		element.Attributes["id"] = generateId(element)
	}

}

func (element *Element) SetClass(value string) {
	element.Attributes["class"] = strings.TrimSpace(value)
}

func (element *Element) SetLine(value int) {
	element.Line = value
}

func (element *Element) SetText(value string) {

	element.Text = strings.TrimSpace(value)

	if element.Type == "h1" || element.Type == "h2" || element.Type == "h3" || element.Type == "h4" {
		element.Attributes["id"] = generateId(element)
	}

}

func (element *Element) stringChildren(space string) string {

	result := ""

	for c := 0; c < len(element.Children); c++ {

		result += element.Children[c].String()

		if c < len(element.Children)-1 {
			result += space
		}

	}

	return result

}

func (element *Element) stringAttributes() string {

	result := ""

	for _, name := range element.getAttributes() {

		value := element.GetAttribute(name)

		if value != "" {
			result += " " + name + "=\"" + value + "\""
		} else {
			result += " " + name
		}

	}

	return strings.TrimSpace(result)

}

func (element *Element) String() string {

	result := ""

	switch element.Type {

	case "a":

		if strings.HasPrefix(element.Text, "#") {
			result = fmt.Sprintf("[%s](%s)", element.Text, element.Text)
		} else {
			result = fmt.Sprintf("[%s](%s)", element.Text, element.GetAttribute("href"))
		}

	case "abbr":

		result = fmt.Sprintf("[%s]{%s}", element.Text, element.GetAttribute("title"))

	case "audio":

		result = fmt.Sprintf("![%s](%s)", element.GetAttribute("title"), element.GetAttribute("src"))

	case "b":

		result = fmt.Sprintf("**%s**", element.Text)

	case "code":

		result = fmt.Sprintf("`%s`", element.Text)

	case "del":

		result = fmt.Sprintf("~%s~", element.Text)

	case "em":

		result = fmt.Sprintf("*%s*", element.Text)

	case "u":

		result = fmt.Sprintf("++%s++", element.Text)

	case "h1":

		result = fmt.Sprintf("# %s", element.stringChildren(" "))

	case "h2":

		result = fmt.Sprintf("## %s", element.stringChildren(" "))

	case "h3":

		result = fmt.Sprintf("### %s", element.stringChildren(" "))

	case "h4":

		result = fmt.Sprintf("#### %s", element.stringChildren(" "))

	case "h5":

		result = fmt.Sprintf("##### %s", element.stringChildren(" "))

	case "h6":

		result = fmt.Sprintf("###### %s", element.stringChildren(" "))

	case "blockquote":

		lines := strings.Split(element.stringChildren("\n"), "\n")

		for index := 0; index < len(lines); index++ {
			lines[index] = "> " + strings.TrimSpace(lines[index])
		}

		result = strings.Join(lines, "\n")

	case "hr":

		result = "---"

	case "img":

		result = fmt.Sprintf("![%s](%s)", element.GetAttribute("alt"), element.GetAttribute("src"))

	case "p":

		content := make([]string, 0)
		current := ""

		for c := 0; c < len(element.Children); c++ {

			if element.Children[c].Type == "#text" {

				words := strings.Split(element.Children[c].String(), " ")

				for w, word := range words {

					current += word

					if len(current) > markdownLineWidth {

						content = append(content, current)
						current = ""

					} else {

						if w < len(words)-1 {
							current += " "
						}

					}

				}

				if current != "" && c < len(element.Children)-1 {
					current += " "
				}

			} else {

				current += element.Children[c].String()

				if len(current) > markdownLineWidth {

					content = append(content, current)
					current = ""

				} else {

					if current != "" && c < len(element.Children)-1 {
						current += " "
					}

				}

			}

		}

		if current != "" {

			content = append(content, current)
			current = ""

		}

		result = fmt.Sprintf("%s", strings.Join(content, "\n"))

	case "pre":

		class, ok := element.Attributes["class"]

		if ok == true {
			result = fmt.Sprintf("```%s\n%s\n```", class, element.Text)
		} else {
			result = fmt.Sprintf("```\n%s\n```", element.Text)
		}

	// Raw HTML nodes
	case "article", "aside", "dialog", "div", "figure", "figcaption", "footer", "header", "main", "span":

		content := element.stringChildren("\n")
		attributes := element.stringAttributes()

		if attributes != "" {
			result = fmt.Sprintf("<%s %s>\n%s\n</%s>", element.Type, attributes, content, element.Type)
		} else {
			result = fmt.Sprintf("<%s>\n%s\n</%s>", element.Type, content, element.Type)
		}

	case "ol":

		content := ""

		for c := 0; c < len(element.Children); c++ {
			content += fmt.Sprintf("%d. %s\n", c+1, element.Children[c].String())
		}

		result = fmt.Sprintf("%s", content)

	case "ul":

		content := ""

		for c := 0; c < len(element.Children); c++ {
			content += fmt.Sprintf("- %s\n", element.Children[c].String())
		}

		result = fmt.Sprintf("%s", content)

	case "li":

		result = fmt.Sprintf("%s", element.stringChildren(" "))

	case "#text":

		result = fmt.Sprintf("%s", element.Text)

	case "table":

		result = stringTable(element)

	default:

		// Do Nothing, unallowed HTML element

	}

	return result

}

func (element *Element) textChildren(space string) string {

	result := ""

	for c := 0; c < len(element.Children); c++ {

		result += element.Children[c].TextContent()

		if c < len(element.Children)-1 {
			result += space
		}

	}

	return result

}

func (element *Element) TextContent() string {

	result := ""

	switch element.Type {

	case "a":

		result = element.Text

	case "abbr":

		result = element.Text

	case "audio":

		result = element.GetAttribute("title")

	case "b":

		result = element.Text

	case "code":

		result = element.Text

	case "del":

		result = element.Text

	case "em":

		result = element.Text

	case "h1", "h2", "h3", "h4", "h5", "h6":

		result = element.textChildren(" ")

	case "u":

		result = element.Text

	case "blockquote":

		result = element.textChildren("\n")

	case "hr":

		result = ""

	case "img":

		result = element.GetAttribute("alt")

	case "p":

		result = element.textChildren(" ")

	case "pre":

		result = element.Text

	case "article", "aside", "dialog", "div", "figure", "figcaption", "footer", "header", "main", "span":

		result = element.textChildren("\n")

	case "ol":

		content := ""

		for c, child := range element.Children {
			content += fmt.Sprintf("%d. %s\n", c+1, child.TextContent())
		}

		result = fmt.Sprintf("%s", content)

	case "ul":

		content := ""

		for _, child := range element.Children {
			content += fmt.Sprintf("- %s\n", child.TextContent())
		}

		result = fmt.Sprintf("%s", content)

	case "li":

		result = element.textChildren(" ")

	case "#text":

		result = element.Text

	case "table":

		result = stringTable(element)

	default:

		// Do Nothing, unallowed HTML element

	}

	return result

}
