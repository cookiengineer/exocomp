package html

import "strings"
import html_tree "golang.org/x/net/html"

func firstHeading(node *html_tree.Node) string {

	if node == nil {
		return ""
	}

	if node.Type == html_tree.ElementNode {

		switch node.Data {
		case "h1", "h2", "h3":
			text := strings.TrimSpace(rawText(node))

			if text != "" {
				return text
			}
		}

	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {

		heading := firstHeading(child)

		if heading != "" {
			return heading
		}

	}

	return ""

}
