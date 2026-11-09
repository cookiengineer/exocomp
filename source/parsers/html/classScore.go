package html

import "strings"

func classScore(class string, id string) float64 {

	combined := strings.ToLower(strings.TrimSpace(class + " " + id))

	if combined == "" {
		return 0
	}

	score := 0.0

	positive := []string{
		"article", "content", "post", "story", "entry", "main", "body",
		"text", "blog", "news", "markdown", "hentry", "page",
	}

	negative := []string{
		"comment", "sidebar", "footer", "header", "nav", "menu", "aside",
		"banner", "advert", "promo", "share", "social", "widget", "related",
		"copyright", "breadcrumb", "cookie", "popup", "modal", "login", "search",
	}

	for _, keyword := range positive {

		if strings.Contains(combined, keyword) {
			score += 25
		}

	}

	for _, keyword := range negative {

		if strings.Contains(combined, keyword) {
			score -= 25
		}

	}

	return score

}
