package html

func isSkippedElement(name string) bool {

	switch name {

	case "script", "style", "noscript", "head", "meta", "link", "base", "title",
		"iframe", "svg", "canvas", "template", "object", "embed", "audio", "video",
		"form", "input", "button", "select", "option", "textarea", "label", "fieldset":
		return true
	default:
		return false

	}

}
