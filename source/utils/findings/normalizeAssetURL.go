package findings

import "net/url"
import "strings"

func normalizeAssetURL(value string) string {

	parsed, err := url.Parse(value)

	if err == nil && parsed.Scheme != "" && parsed.Host != "" {

		parsed.Scheme = strings.ToLower(parsed.Scheme)
		parsed.Host = strings.ToLower(parsed.Host)

		return parsed.String()

	} else {
		return value
	}

}
