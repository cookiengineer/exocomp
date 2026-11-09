package html

import html_tree "golang.org/x/net/html"

func findBodyNode(node *html_tree.Node) *html_tree.Node {

	if node == nil {
		return nil
	}

	if node.Type == html_tree.ElementNode && node.Data == "body" {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {

		found := findBodyNode(child)

		if found != nil {
			return found
		}

	}

	return nil

}
