package html

func isBlockElement(name string) bool {

	switch name {

	case "h1", "h2", "h3", "h4", "h5", "h6", "p", "ul", "ol", "li",
		"pre", "blockquote", "table", "tr", "td", "th", "hr":
		return true
	default:
		return false

	}

}
