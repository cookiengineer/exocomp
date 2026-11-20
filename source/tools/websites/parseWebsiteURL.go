package websites

import "fmt"
import net_url "net/url"
import "strings"

func parseWebsiteURL(raw string) (*net_url.URL, error) {

	url, err := net_url.Parse(strings.TrimSpace(raw))

	if err == nil {

		if url.Scheme == "http" || url.Scheme == "https" {

			if url.Host != "" {
				return url, nil
			} else {
				return nil, fmt.Errorf("websites.ParseURL: Invalid URL \"%s\": Missing host.", raw)
			}
		} else {
			return nil, fmt.Errorf("websites.ParseURL: Invalid URL \"%s\": Scheme must be http or https.", raw)
		}
	} else {
		return nil, fmt.Errorf("websites.ParseURL: Invalid URL \"%s\".", raw)
	}
}
