package html

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
