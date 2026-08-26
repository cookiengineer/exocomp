package tools

import "fmt"
import net_url "net/url"
import "strings"

func parseWebsiteURL(raw string) (*net_url.URL, error) {

	url, err := net_url.Parse(strings.TrimSpace(raw))

	if err != nil {
		return nil, fmt.Errorf("websites: Invalid URL \"%s\"", raw)
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("websites: Invalid URL \"%s\": Scheme must be http or https", raw)
	}

	if url.Host == "" {
		return nil, fmt.Errorf("websites: Invalid URL \"%s\": Missing host", raw)
	}

	return url, nil

}

