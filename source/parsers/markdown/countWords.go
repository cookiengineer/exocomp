package markdown

import "strings"

func countWords(element *Element) int {

	result := 0

	if element.Text != "" {

		words := strings.Split(element.Text, " ")

		for w := 0; w < len(words); w++ {
			result++
		}

	}

	if len(element.Children) > 0 {

		for c := 0; c < len(element.Children); c++ {
			result += countWords(element.Children[c])
		}

	}

	return result

}
