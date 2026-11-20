package assets

import "net/url"
import "strings"

func hostFromValue(value string) string {

	if strings.Contains(value, "://") == true {

		if parsed, err := url.Parse(value); err == nil {
			return parsed.Hostname()
		}

	}

	if strings.Contains(value, "/") == true {
		return strings.Split(value, "/")[0]
	} else {
		return value
	}

}
