package html

import "strings"
import html_tree "golang.org/x/net/html"

func getAttribute(node *html_tree.Node, key string) string {

	for _, attribute := range node.Attr {

		if strings.EqualFold(attribute.Key, key) {
			return attribute.Val
		}

	}

	return ""

}
