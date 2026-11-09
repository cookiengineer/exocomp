package html

import "strings"
import html_tree "golang.org/x/net/html"

func rawText(node *html_tree.Node) string {

	if node == nil {
		return ""
	}

	var result strings.Builder

	var walk func(node *html_tree.Node)

	walk = func(node *html_tree.Node) {

		if node.Type == html_tree.TextNode {
			result.WriteString(node.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}

	}

	walk(node)

	return result.String()

}
