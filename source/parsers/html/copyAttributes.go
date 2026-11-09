package html

import "strings"
import html_tree "golang.org/x/net/html"

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
