package tools

import parser_html "exocomp/parsers/html"
import net_url "net/url"
import "strings"

func getWebsiteText(url *net_url.URL, bytes []byte) string {

	document, err := parser_html.Parse(url.String(), bytes)

	if err == nil && document != nil {
		return document.Text()
	}

	return strings.TrimSpace(string(bytes))

}
