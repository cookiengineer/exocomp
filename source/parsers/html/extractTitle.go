package html

import "strings"
import html_tree "golang.org/x/net/html"

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
