package markdown

import "fmt"
import net_url "net/url"
import "regexp"
import "sort"
import "strings"
import "time"

type Document struct {
	URL  *net_url.URL `json:"url"`
	Meta struct {
		Title   string    `json:"title"`
		Summary string    `json:"summary"`
		Date    time.Time `json:"date"`
		Tags    []string  `json:"tags"`
		Image   string    `json:"image"`
	} `json:"meta"`
	Statistics struct {
		Minutes int `json:"minutes"`
		Words   int `json:"words"`
	} `json:"statistics"`
	Body   []*Element `json:"body"`
	Errors []error    `json:"errors"`
}

func NewDocument(file string) *Document {

	base := strings.TrimSpace(file)

	if base == "" {
		base = "https://cookie.engineer/index.html"
	} else if strings.HasPrefix(base, "/") {
		base = "https://cookie.engineer/" + base[1:]
	}

	base_url, err := net_url.Parse(base)

	if err == nil {

		var document Document

		document.URL = base_url
		document.Body = make([]*Element, 0)
		document.Errors = make([]error, 0)
		document.Meta.Image = "https://cookie.engineer/design/about/cookiengineer.jpg"
		document.Meta.Tags = make([]string, 0)

		return &document

	}

	return nil

}

func (document *Document) AddElement(element *Element) {

	if element != nil && element.Type != "" {
		document.Body = append(document.Body, element)
	}

}

func (document *Document) Count() {

	var words int = 0

	for b := 0; b < len(document.Body); b++ {
		words += countWords(document.Body[b])
	}

	document.Statistics.Words = words
	document.Statistics.Minutes = int(document.Statistics.Words / 200)

}

func (document *Document) CloseElement() bool {

	var result bool = false
	var element *Element = nil

	if len(document.Body) > 0 {
		element = document.Body[len(document.Body)-1]
	}

	if element != nil {

		if element.Type != "#text" {

			text := NewElement("#text")
			text.SetLine(element.Line + 1)

			document.Body = append(document.Body, text)

		}

	}

	return result

}

func (document *Document) getLastElement() *Element {

	var element *Element = nil

	if len(document.Body) > 0 {
		element = document.Body[len(document.Body)-1]
	}

	return element

}

func (document *Document) IsValid() bool {

	var result bool

	if document.URL != nil && len(document.Body) > 0 {

		if document.Meta.Title != "" && document.Meta.Summary != "" {

			if document.Meta.Date.IsZero() == false && len(document.Meta.Tags) > 0 {

				if len(document.Errors) == 0 {
					result = true
				}

			}

		}

	}

	return result

}

func (document *Document) MarshalText() ([]byte, error) {

	return []byte(document.String()), nil

}

func (document *Document) Parse(bytes []byte) []error {

	errors := make([]error, 0)

	markdown := strings.TrimSpace(string(bytes))

	if markdown != "" {

		errors_meta := document.ParseMeta(markdown)

		for _, err := range errors_meta {
			errors = append(errors, err)
		}

		errors_body := document.ParseBody(markdown)

		for _, err := range errors_body {
			errors = append(errors, err)
		}

	} else {
		errors = append(errors, fmt.Errorf("Line 1: Expected non-empty markdown code"))
	}

	return errors

}

func (document *Document) ParseMeta(value string) []error {

	errors := make([]error, 0)
	lines := strings.Split(strings.TrimSpace(value), "\n")
	found := map[string]bool{
		"title":   false,
		"summary": false,
		"date":    false,
		"tags":    false,
	}

	if lines[0] == "===" {

		for l := 1; l < len(lines); l++ {

			line := strings.TrimSpace(lines[l])

			if strings.HasPrefix(line, "- ") && strings.Contains(line, ":") {

				key := strings.TrimSpace(line[2:strings.Index(line, ":")])
				val := strings.TrimSpace(line[strings.Index(line, ":")+1:])

				if key == "title" {

					document.SetTitle(val)
					found["title"] = true

				} else if key == "summary" {

					document.SetSummary(val)
					found["summary"] = true

				} else if key == "date" {

					date, err := time.Parse("2006-01-02", val)

					if err == nil {
						document.SetDate(date)
						found["date"] = true
					}

				} else if key == "tags" {

					values := strings.Split(val, ",")

					if len(values) > 0 && values[0] != "" {
						document.SetTags(values)
						found["tags"] = true
					}

				} else if key == "image" {

					document.SetImage(val)

				}

			} else if line == "===" {
				break
			}

		}

		if found["title"] == false {
			errors = append(errors, fmt.Errorf("Line 1: Missing \"title\" (string) frontmatter field"))
		}

		if found["summary"] == false {
			errors = append(errors, fmt.Errorf("Line 1: Missing \"summary\" (string) frontmatter field"))
		}

		if found["date"] == false {
			errors = append(errors, fmt.Errorf("Line 1: Missing \"date\" (YYYY-MM-DD string) frontmatter field"))
		}

		if found["tags"] == false {
			errors = append(errors, fmt.Errorf("Line 1: Missing \"tags\" (comma separated strings) frontmatter field"))
		}

	}

	return errors

}

func (document *Document) ParseBody(value string) []error {

	errors := make([]error, 0)
	lines := strings.Split(strings.TrimSpace(value), "\n")
	regexp_ul, _ := regexp.Compile("^([\\*\\-+]+)[[:space:]]")
	regexp_ol, _ := regexp.Compile("^([0-9]+)\\.[[:space:]]")

	first_line := 0

	if lines[0] == "===" {

		for l := 1; l < len(lines); l++ {

			line := strings.TrimSpace(lines[l])

			if line == "===" {
				first_line = l + 1
				break
			}

		}

	}

	for l := first_line; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])
		pointer := document.getLastElement()

		if pointer != nil && pointer.Type == "pre" {

			if line != "```" {
				// Preserve indention
				pointer.AddText(lines[l])
			} else {
				document.CloseElement()
			}

		} else if pointer != nil && pointer.Type == "table" {

			if strings.HasPrefix(line, "|") {
				pointer.AddText(line)
			} else {
				document.CloseElement()
			}

		} else if line == "" {

			if pointer == nil {
				// Do Nothing
			} else if pointer.Type == "pre" {
				// Do Nothing
			} else if pointer.Type == "div" || pointer.Type == "figure" || pointer.Type == "figcaption" || pointer.Type == "p" {
				document.CloseElement()
			} else if pointer.Type == "article" || pointer.Type == "dialog" || pointer.Type == "footer" || pointer.Type == "header" || pointer.Type == "main" {
				document.CloseElement()
			} else if pointer.Type == "ul" || pointer.Type == "ol" {
				document.CloseElement()
			} else if pointer.Type == "table" {
				document.CloseElement()
			} else if pointer != nil {
				document.CloseElement()
				pointer = nil
			}

		} else if strings.HasPrefix(line, "![") && strings.Contains(line, "](") && strings.HasSuffix(line, ")") {

			text := line[2:strings.Index(line, "](")]
			href := line[strings.Index(line, "](")+2 : strings.Index(line, ")")]

			if strings.HasSuffix(href, ".m4a") || strings.HasSuffix(href, ".mp3") || strings.HasSuffix(href, ".opus") {

				element := NewElement("audio")
				element.SetLine(l + 1)
				element.SetAttribute("controls", "")
				element.SetAttribute("src", href)
				element.SetAttribute("title", text)

				if pointer != nil && pointer.IsBlockElement() == true {
					pointer.AddChild(element)
				} else {
					document.AddElement(element)
				}

			} else if strings.HasSuffix(href, ".gif") || strings.HasSuffix(href, ".jpg") || strings.HasSuffix(href, ".png") {

				element := NewElement("img")
				element.SetLine(l + 1)
				element.SetAttribute("alt", text)
				element.SetAttribute("src", href)

				if pointer != nil && pointer.IsBlockElement() == true {
					pointer.AddChild(element)
				} else {
					document.AddElement(element)
				}

			}

		} else if strings.HasPrefix(line, "[") && strings.Contains(line, "]: ") {

			name := strings.TrimSpace(line[1:strings.Index(line, "]: ")])
			children, err := parseInlineElements(line[strings.Index(line, "]: ")+3:], l+1, true)

			if err == nil {

				if pointer == nil {

					anchor := NewElement("span")
					anchor.SetLine(l + 1)
					anchor.SetText(fmt.Sprintf("[%s]:", name))
					anchor.SetAttribute("id", fmt.Sprintf("footnote-%s", name))

					element := NewElement("p")
					element.SetLine(l + 1)
					element.AddChild(anchor)
					element.AddChildren(children)

					document.AddElement(element)

				} else if pointer.Type == "p" {

					anchor := NewElement("span")
					anchor.SetLine(l + 1)
					anchor.SetText(fmt.Sprintf("[%s]:", name))
					anchor.SetAttribute("id", fmt.Sprintf("footnote-%s", name))

					pointer.AddChild(anchor)
					pointer.AddChildren(children)

				} else if pointer != nil && pointer.IsBlockElement() == true {

					anchor := NewElement("span")
					anchor.SetLine(l + 1)
					anchor.SetText(fmt.Sprintf("[%s]:", name))
					anchor.SetAttribute("id", fmt.Sprintf("footnote-%s", name))

					pointer.AddChild(anchor)
					pointer.AddChildren(children)

				} else {

					anchor := NewElement("span")
					anchor.SetLine(l + 1)
					anchor.SetText(fmt.Sprintf("[%s]:", name))
					anchor.SetAttribute("id", fmt.Sprintf("footnote-%s", name))

					element := NewElement("p")
					element.SetLine(l + 1)
					element.AddChild(anchor)
					element.AddChildren(children)

					document.AddElement(element)

				}

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if len(line) >= 3 && strings.HasPrefix(line, "```") {

			if pointer == nil || pointer.Type != "pre" {

				class := strings.TrimSpace(line[3:])

				if class == "" {
					errors = append(errors, fmt.Errorf("Line %d: Expected non-empty language attribute", l+1))
				}

				element := NewElement("pre")
				element.SetLine(l + 1)
				element.SetAttribute("class", class)

				document.AddElement(element)

			}

		} else if len(line) > 3 && strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {

			if pointer == nil || pointer.Type != "table" {

				element := NewElement("table")
				element.SetLine(l + 1)
				element.SetText(line)

				document.AddElement(element)

			}

		} else if strings.HasPrefix(line, "#####") {

			children, err := parseInlineElements(strings.TrimSpace(line[5:]), l+1, false)

			if err == nil {

				element := NewElement("h5")
				element.SetLine(l + 1)
				element.SetChildren(children)

				document.AddElement(element)

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if strings.HasPrefix(line, "####") {

			children, err := parseInlineElements(strings.TrimSpace(line[4:]), l+1, false)

			if err == nil {

				element := NewElement("h4")
				element.SetLine(l + 1)
				element.SetChildren(children)

				document.AddElement(element)

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if strings.HasPrefix(line, "###") {

			children, err := parseInlineElements(strings.TrimSpace(line[3:]), l+1, false)

			if err == nil {

				element := NewElement("h3")
				element.SetLine(l + 1)
				element.SetChildren(children)

				document.AddElement(element)

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if strings.HasPrefix(line, "##") {

			children, err := parseInlineElements(strings.TrimSpace(line[2:]), l+1, false)

			if err == nil {

				element := NewElement("h2")
				element.SetLine(l + 1)
				element.SetChildren(children)

				document.AddElement(element)

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if strings.HasPrefix(line, "#") {

			children, err := parseInlineElements(strings.TrimSpace(line[1:]), l+1, false)

			if err == nil {

				element := NewElement("h1")
				element.SetLine(l + 1)
				element.SetChildren(children)

				document.AddElement(element)

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {

			if strings.HasPrefix(line, "</") && strings.HasSuffix(line, ">") {

				element_type := strings.TrimSpace(line[2 : len(line)-1])

				if element_type != "" {

					if pointer != nil && pointer.Type == element_type {
						pointer = nil
					}

				} else {
					errors = append(errors, fmt.Errorf("Line %d: Malformed raw HTML tag \"%s\"", l+1, line))
				}

			} else if strings.HasPrefix(line, "<") && strings.HasSuffix(line, "/>") {

				element_type := ""

				if strings.Contains(line, " ") {
					element_type = strings.TrimSpace(line[1:strings.Index(line, " ")])
				} else {
					element_type = strings.TrimSpace(line[1 : len(line)-2])
				}

				if element_type != "" {

					element := NewElement(element_type)
					element.SetLine(l + 1)

					// Don't expect a closing HTML tag
					element.is_block_element = false

					if strings.Contains(line, " ") {

						raw_attributes := strings.Split(line[strings.Index(line, " ")+1:len(line)-2], " ")

						for _, attribute := range raw_attributes {

							if strings.Contains(attribute, "=") {

								attribute_key := attribute[0:strings.Index(attribute, "=")]
								attribute_val := attribute[strings.Index(attribute, "=")+1:]

								if strings.HasPrefix(attribute_val, "\"") && strings.HasSuffix(attribute_val, "\"") {
									attribute_val = attribute_val[1 : len(attribute_val)-1]
								} else if strings.HasPrefix(attribute_val, "'") && strings.HasSuffix(attribute_val, "'") {
									attribute_val = attribute_val[1 : len(attribute_val)-1]
								}

								element.SetAttribute(attribute_key, attribute_val)

							} else if strings.TrimSpace(attribute) != "" {

								element.SetAttribute(strings.TrimSpace(attribute), "")

							}

						}

					}

					document.AddElement(element)

				} else {
					errors = append(errors, fmt.Errorf("Line %d: Malformed raw HTML tag \"%s\"", l+1, line))
				}

			} else if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {

				element_type := ""

				if strings.Contains(line, " ") {
					element_type = line[1:strings.Index(line, " ")]
				} else {
					element_type = line[1 : len(line)-1]
				}

				if element_type != "" {

					element := NewElement(element_type)
					element.SetLine(l + 1)

					// Expect a closing HTML tag
					element.is_block_element = true

					if strings.Contains(line, " ") {

						raw_attributes := strings.Split(line[strings.Index(line, " ")+1:len(line)-1], " ")

						for _, attribute := range raw_attributes {

							if strings.Contains(attribute, "=") {

								attribute_key := attribute[0:strings.Index(attribute, "=")]
								attribute_val := attribute[strings.Index(attribute, "=")+1:]

								if strings.HasPrefix(attribute_val, "\"") && strings.HasSuffix(attribute_val, "\"") {
									attribute_val = attribute_val[1 : len(attribute_val)-1]
								} else if strings.HasPrefix(attribute_val, "'") && strings.HasSuffix(attribute_val, "'") {
									attribute_val = attribute_val[1 : len(attribute_val)-1]
								}

								element.SetAttribute(attribute_key, attribute_val)

							} else if strings.TrimSpace(attribute) != "" {
								element.SetAttribute(strings.TrimSpace(attribute), "")
							}

						}

					}

					document.AddElement(element)

					pointer = document.getLastElement()

				} else {
					errors = append(errors, fmt.Errorf("Line %d: Malformed raw HTML tag \"%s\"", l+1, line))
				}

			}

		} else if regexp_ul.MatchString(line) {

			children, err := parseInlineElements(strings.TrimSpace(line[2:]), l+1, true)

			if err == nil {

				item := NewElement("li")
				item.SetLine(l + 1)
				item.SetChildren(children)

				if pointer == nil {

					element := NewElement("ul")
					element.SetLine(l + 1)

					document.AddElement(element)

					pointer = document.getLastElement()
					pointer.AddChild(item)

				} else if pointer.Type == "ul" {

					pointer.AddChild(item)

				} else if pointer != nil && pointer.IsBlockElement() == true {

					inline_pointer := pointer.getLastChild()

					if inline_pointer == nil || inline_pointer.Type != "ul" {

						element := NewElement("ul")
						element.SetLine(l + 1)

						pointer.AddChild(element)

						inline_pointer = pointer.getLastChild()

					}

					inline_pointer.AddChild(item)

				} else {

					element := NewElement("ul")
					element.SetLine(l + 1)

					document.AddElement(element)

					pointer = document.getLastElement()
					pointer.AddChild(item)

				}

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if regexp_ol.MatchString(line) {

			children, err := parseInlineElements(strings.TrimSpace(line[strings.Index(line, ".")+1:]), l+1, true)

			if err == nil {

				item := NewElement("li")
				item.SetLine(l + 1)
				item.SetChildren(children)

				if pointer == nil {

					element := NewElement("ol")
					element.SetLine(l + 1)

					document.AddElement(element)

					pointer = document.getLastElement()
					pointer.AddChild(item)

				} else if pointer.Type == "ol" {

					pointer.AddChild(item)

				} else if pointer != nil && pointer.IsBlockElement() == true {

					inline_pointer := pointer.getLastChild()

					if inline_pointer == nil || inline_pointer.Type != "ol" {

						element := NewElement("ol")
						element.SetLine(l + 1)

						pointer.AddChild(element)

						inline_pointer = pointer.getLastChild()

					}

					inline_pointer.AddChild(item)

				} else {

					element := NewElement("ol")
					element.SetLine(l + 1)

					document.AddElement(element)

					pointer = document.getLastElement()
					pointer.AddChild(item)

				}

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		} else if line != "" {

			children, err := parseInlineElements(strings.TrimSpace(line), l+1, true)

			if err == nil {

				if pointer == nil {

					element := NewElement("p")
					element.SetLine(l + 1)
					element.AddChildren(children)

					document.AddElement(element)

				} else if pointer.Type == "p" {

					pointer.AddChildren(children)

				} else if pointer != nil && pointer.IsBlockElement() == true {

					pointer.AddChildren(children)

				} else {

					element := NewElement("p")
					element.SetLine(l + 1)
					element.AddChildren(children)

					document.AddElement(element)

				}

			} else {
				errors = append(errors, fmt.Errorf("Line %d: %s", l+1, err.Error()))
			}

		}

	}

	return errors

}

func (document *Document) SetDate(value time.Time) bool {

	if value.IsZero() == false {

		document.Meta.Date = value

		return true

	}

	return false

}

func (document *Document) SetImage(value string) bool {

	tmp := strings.TrimSpace(value)

	if strings.HasPrefix(value, "./") || strings.HasPrefix(value, "/") {

		document.Meta.Image = tmp

		return true

	}

	return false

}

func (document *Document) SetSummary(value string) bool {

	tmp := strings.TrimSpace(value)

	if tmp != "" {

		document.Meta.Summary = tmp

		return true

	}

	return false

}

func (document *Document) SetTags(values []string) {

	filtered := make([]string, 0)

	for _, value := range values {

		tmp := strings.TrimSpace(strings.ToLower(value))

		if tmp != "" {
			filtered = append(filtered, tmp)
		}

	}

	document.Meta.Tags = filtered

}

func (document *Document) SetTitle(value string) {
	document.Meta.Title = strings.TrimSpace(value)
}

func (document *Document) Render(indent string) string {

	result := make([]string, 0)

	first_section_at := 0

	for e, element := range document.Body {

		if element.Type == "h3" {
			first_section_at = e
			break
		}

	}

	for e, element := range document.Body {

		if element.Type == "h1" {
			result = append(result, element.RenderInto(document, indent+"\t"))
		} else if element.Type == "h2" {
			result = append(result, element.RenderInto(document, indent+"\t"))
		} else if element.Type == "h3" {

			if first_section_at == e {
				result = append(result, indent+"\t<section>")
			} else {
				result = append(result, indent+"\t</section>")
				result = append(result, indent+"\t<section>")
			}

			result = append(result, element.RenderInto(document, indent+"\t\t"))

		} else if element.Type == "h4" {
			result = append(result, element.RenderInto(document, indent+"\t\t"))
		} else if element.Type == "h5" {
			result = append(result, element.RenderInto(document, indent+"\t\t"))
		} else if element.Type != "#text" {

			if e < first_section_at {
				result = append(result, element.RenderInto(document, indent+"\t"))
			} else {
				result = append(result, element.RenderInto(document, indent+"\t\t"))
			}

		}

	}

	result = append(result, indent+"\t</section>")

	return strings.Join(result, "\n")

}

func (document *Document) String() string {

	lines := make([]string, 0)

	lines = append(lines, "===")

	if document.Meta.Title != "" {
		lines = append(lines, fmt.Sprintf("- title: %s", document.Meta.Title))
	}

	if document.Meta.Summary != "" {
		lines = append(lines, fmt.Sprintf("- summary: %s", document.Meta.Summary))
	}

	if document.Meta.Date.IsZero() == false {
		lines = append(lines, fmt.Sprintf("- date: %s", document.Meta.Date.Format("2006-01-02")))
	}

	if len(document.Meta.Tags) > 0 {

		tags := make([]string, 0)

		for _, tag := range document.Meta.Tags {
			tags = append(tags, tag)
		}

		sort.Strings(tags)

		lines = append(lines, fmt.Sprintf("- tags: %s", strings.Join(tags, ", ")))

	}

	if document.Meta.Image != "" {
		lines = append(lines, fmt.Sprintf("- image: %s", document.Meta.Image))
	}

	lines = append(lines, "===")
	lines = append(lines, "")

	for _, element := range document.Body {

		switch element.Type {

		case "h1", "h2", "h3", "h4", "h5":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "audio", "img":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "p":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "pre":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "article", "aside", "dialog", "div", "figure", "figcaption", "footer", "header", "main":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "span":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			lines = append(lines, "")

		case "ol", "ul":

			lines = append(lines, fmt.Sprintf("%s", element.String()))
			// XXX: Unsure if this is needed or not
			// lines = append(lines, "")

		case "table":

			lines = append(lines, fmt.Sprintf("%s", strings.TrimSpace(element.String())))
			lines = append(lines, "")

		}

	}

	return strings.Join(lines, "\n")

}

func (document *Document) UnmarshalMarkdown(bytes []byte) error {

	// Set default image
	document.SetImage("/design/about/cookiengineer.jpg")

	document.Errors = document.Parse(bytes)

	if len(document.Errors) > 0 {

		messages := make([]string, 0, len(document.Errors))

		for _, err := range document.Errors {
			messages = append(messages, err.Error())
		}

		return fmt.Errorf("%d Parsing Errors:\n%s\n", len(messages), strings.Join(messages, "\n"))

	} else {

		return nil

	}

}
