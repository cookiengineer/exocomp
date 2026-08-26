package tools

import parser_html "exocomp/parsers/html"
import net_url "net/url"
import "strings"

func getWebsiteMarkdown(url *net_url.URL, bytes []byte) string {

	document, err := parser_html.Parse(url.String(), bytes)

	if err == nil && document != nil {
		return document.String()
	}

	return strings.TrimSpace(string(bytes))

}
