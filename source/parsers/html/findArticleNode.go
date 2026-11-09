package html

import "math"
import html_tree "golang.org/x/net/html"

// findArticleNode locates the node that most likely contains the main
// article content, using a Readability-style scoring heuristic. Returns nil
// when no meaningful content is found.
func findArticleNode(root *html_tree.Node) *html_tree.Node {

	if root == nil {
		return nil
	}

	scores := make(map[*html_tree.Node]float64)
	depths := make(map[*html_tree.Node]int)

	var score func(node *html_tree.Node, depth int) float64

	score = func(node *html_tree.Node, depth int) float64 {

		own := 0.0

		if node.Type == html_tree.ElementNode {

			text := rawText(node)

			switch node.Data {

			case "p":
				own = 1 + commaScore(text) + math.Min(float64(len(text))/100.0, 3.0)
			case "pre", "blockquote":
				own = 1 + float64(len(text))/200.0
			case "h1", "h2", "h3", "h4", "h5", "h6":
				own = 1
			case "li", "td", "th":
				own = 0.5
			case "article", "main":
				own = 10
			case "nav", "header", "footer", "aside":
				own = -5

			}

			own += classScore(getAttribute(node, "class"), getAttribute(node, "id"))

		}

		total := own

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			total += score(child, depth+1)
		}

		scores[node] = total
		depths[node] = depth

		return total

	}

	score(root, 0)

	var best *html_tree.Node
	var best_score float64 = 0

	for node, value := range scores {

		if node.Type != html_tree.ElementNode {
			continue
		}

		switch node.Data {
		case "html", "head":
			continue
		}

		if best == nil || value > best_score || (value == best_score && depths[node] > depths[best]) {

			best = node
			best_score = value

		}

	}

	return best

}
