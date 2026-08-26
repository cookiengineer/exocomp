package types

import "net/http"
import "strings"

type UserAgent struct {
	Name      string            `json:"name" yaml:"name"`
	UserAgent string            `json:"user-agent" yaml:"user-agent"`
	Platform  string            `json:"platform" yaml:"platform"`
	Mobile    bool              `json:"mobile" yaml:"mobile"`
	Headers   map[string]string `json:"headers" yaml:"headers"`
}

func NewUserAgent(name string, user_agent string, platform string, mobile bool, headers map[string]string) *UserAgent {

	useragent := &UserAgent{
		Name:      strings.TrimSpace(name),
		UserAgent: strings.TrimSpace(user_agent),
		Platform:  strings.TrimSpace(platform),
		Mobile:    mobile,
		Headers:   make(map[string]string),
	}

	for key, value := range headers {
		useragent.Headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return useragent

}

// Header returns the HTTP headers a request should send to impersonate this
// browser, combining the User-Agent string with its Client Hints (SEC-CH-*).
func (useragent *UserAgent) Header() http.Header {

	header := http.Header{}

	if useragent != nil {

		if useragent.UserAgent != "" {
			header.Set("User-Agent", useragent.UserAgent)
		}

		for name, value := range useragent.Headers {

			if name != "" && value != "" {
				header.Set(name, value)
			}

		}

	}

	return header

}

