package html

import net_url "net/url"
import "strings"

func resolveURL(document *Document, ref string) string {

	ref = strings.TrimSpace(ref)

	if ref == "" {
		return ref
	}

	if document == nil || document.URL == nil {
		return ref
	}

	ref_url, err := net_url.Parse(ref)

	if err != nil {
		return ref
	}

	return document.URL.ResolveReference(ref_url).String()

}
