package html

import net_url "net/url"
import "fmt"
import "strings"

import html_tree "golang.org/x/net/html"

type Document struct {
	URL    *net_url.URL `json:"url"`
	Title  string       `json:"title"`
	Body   []*Element   `json:"body"`
	Errors []error      `json:"errors"`
}

func NewDocument(file string) *Document {

	base := strings.TrimSpace(file)

	if base == "" {
		return nil
	}

	base_url, err := net_url.Parse(base)

	if err == nil {

		document := &Document{
			URL:    base_url,
			Title:  "",
			Body:   make([]*Element, 0),
			Errors: make([]error, 0),
		}

		return document

	}

	return nil

}

func (document *Document) AddElement(element *Element) {

	if document != nil && element != nil && element.Type != "" {
		document.Body = append(document.Body, element)
	}

}

func (document *Document) MarshalText() ([]byte, error) {

	if document == nil {
		return []byte{}, fmt.Errorf("nil document")
	}

	return []byte(document.String()), nil

}

// Parse converts raw HTML bytes into a normalized Document tree, dropping
// script/style/meta content and everything else that is irrelevant for
// reading a page.
func (document *Document) Parse(bytes []byte) error {

	if document == nil {
		return fmt.Errorf("nil document")
	}

	reader  := strings.NewReader(string(bytes))
	node, err := html_tree.Parse(reader)

	if err != nil {
		return err
	}

	document.Errors = make([]error, 0)
	document.Body   = make([]*Element, 0)

	article := findArticleNode(node)

	if article == nil {
		article = findBodyNode(node)
	}

	document.Title = extractTitle(node)

	if document.Title == "" && article != nil {
		document.Title = firstHeading(article)
	}

	document.buildBody(article)

	return nil

}

// buildBody adds the extracted article subtree to the document body. When the
// article root is a container (body/div/article/...) its children are
// flattened into the body; when it is a leaf block (pre/table/h1/...) that
// block becomes a top-level element.
func (document *Document) buildBody(article *html_tree.Node) {

	if article == nil {
		return
	}

	name := strings.ToLower(strings.TrimSpace(article.Data))

	if isTransparentElement(name) {
		document.walk(article, nil)
		return
	}

	if isBlockElement(name) {

		element := NewElement(name)
		copyAttributes(article, element)
		document.AddElement(element)
		document.walk(article, element)
		return

	}

	document.walk(article, nil)

}

func (document *Document) String() string {

	if document == nil {
		return ""
	}

	lines := make([]string, 0)

	if document.Title != "" {
		lines = append(lines, "# "+strings.TrimSpace(document.Title))
		lines = append(lines, "")
	}

	for _, element := range document.Body {

		content := element.Render(document)

		if strings.TrimSpace(content) != "" {
			lines = append(lines, content)
			lines = append(lines, "")
		}

	}

	return strings.TrimSpace(strings.Join(lines, "\n"))

}

// Text returns the plain text content of the document, dropping all markup.
func (document *Document) Text() string {

	if document == nil {
		return ""
	}

	lines := make([]string, 0)

	for _, element := range document.Body {

		text := strings.TrimSpace(element.textContent())

		if text != "" {
			lines = append(lines, text)
		}

	}

	return strings.Join(lines, "\n")

}

// walk traverses the parsed HTML tree and builds a simplified Document tree.
func (document *Document) walk(node *html_tree.Node, current *Element) {

	if node == nil {
		return
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {

		if child.Type == html_tree.TextNode {

			text := normalizeText(child.Data)

			if text != "" {

				if current == nil {
					current = document.newParagraph()
				}

				current.AddChild(NewTextElement(text))

			}

		} else if child.Type == html_tree.ElementNode {

			name := strings.ToLower(strings.TrimSpace(child.Data))

			if name == "" || isSkippedElement(name) {
				continue
			}

			if isTransparentElement(name) {
				document.walk(child, current)
				continue
			}

			element := NewElement(name)
			copyAttributes(child, element)

			if isBlockElement(name) {

				if current == nil {
					document.AddElement(element)
				} else {
					current.AddChild(element)
				}

				document.walk(child, element)

			} else {

				if current == nil {
					current = document.newParagraph()
				}

				current.AddChild(element)
				document.walk(child, element)

			}

		}

	}

}

func (document *Document) newParagraph() *Element {

	element := NewElement("p")
	document.AddElement(element)

	return element

}

func extractTitle(node *html_tree.Node) string {

	if node == nil {
		return ""
	}

	if node.Type == html_tree.ElementNode && strings.EqualFold(node.Data, "title") {
		return strings.TrimSpace(node.FirstChild.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {

		title := extractTitle(child)

		if title != "" {
			return title
		}

	}

	return ""

}

// normalizeText collapses runs of whitespace into a single space, keeping a
// single leading/trailing space when the original text had one, so that
// inline elements stay separated ("Hello <b>world</b>" -> "Hello **world**").
func normalizeText(raw string) string {

	collapsed := strings.Join(strings.Fields(raw), " ")

	if collapsed == "" {
		return ""
	}

	result := collapsed

	if len(raw) > 0 && isWhitespace(raw[0]) {
		result = " " + result
	}

	if len(raw) > 0 && isWhitespace(raw[len(raw)-1]) {
		result = result + " "
	}

	return result

}

func isWhitespace(chr byte) bool {
	return chr == ' ' || chr == '\t' || chr == '\n' || chr == '\r' || chr == '\f' || chr == '\v'
}

func copyAttributes(source *html_tree.Node, target *Element) {
	for _, attribute := range source.Attr {

		key := strings.ToLower(strings.TrimSpace(attribute.Key))
		val := strings.TrimSpace(attribute.Val)

		if key == "class" || key == "id" || key == "style" || key == "onclick" || key == "onload" || key == "onerror" {
			continue
		}

		if key != "" && val != "" {
			target.SetAttribute(key, val)
		}

	}

}

func isSkippedElement(name string) bool {

	switch name {
	case "script", "style", "noscript", "head", "meta", "link", "base", "title",
		"iframe", "svg", "canvas", "template", "object", "embed", "audio", "video",
		"form", "input", "button", "select", "option", "textarea", "label", "fieldset":
		return true
	default:
		return false
	}

}

func isBlockElement(name string) bool {

	switch name {
	case "h1", "h2", "h3", "h4", "h5", "h6", "p", "ul", "ol", "li",
		"pre", "blockquote", "table", "tr", "td", "th", "hr":
		return true
	default:
		return false
	}

}

func isTransparentElement(name string) bool {

	switch name {
	case "html", "body", "div", "section", "article", "main", "aside", "header", "footer", "nav",
		"figure", "figcaption", "address", "details", "summary", "dialog",
		"thead", "tbody", "tfoot", "colgroup", "caption", "dl", "dt", "dd":
		return true
	default:
		return false
	}

}
