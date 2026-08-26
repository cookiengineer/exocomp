package types

import "strings"

var UserAgents []*UserAgent

func GetUserAgent(name string) *UserAgent {

	search := strings.TrimSpace(strings.ToLower(name))

	for _, useragent := range UserAgents {

		if useragent != nil && strings.ToLower(useragent.Name) == search {
			return useragent
		}

	}

	return nil

}

func init() {

	UserAgents = make([]*UserAgent, 0)

	UserAgents = append(UserAgents, NewUserAgent(
		"chrome-windows",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Windows",
		false,
		map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.9",
			"Sec-CH-UA":                 `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			"Sec-CH-UA-Mobile":          "?0",
			"Sec-CH-UA-Platform":        `"Windows"`,
			"Sec-CH-UA-Platform-Version": `"15.0.0"`,
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"chrome-macos",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"macOS",
		false,
		map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.9",
			"Sec-CH-UA":                 `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			"Sec-CH-UA-Mobile":          "?0",
			"Sec-CH-UA-Platform":        `"macOS"`,
			"Sec-CH-UA-Platform-Version": `"14.3.1"`,
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"chrome-linux",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Linux",
		false,
		map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.9",
			"Sec-CH-UA":                 `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			"Sec-CH-UA-Mobile":          "?0",
			"Sec-CH-UA-Platform":        `"Linux"`,
			"Sec-CH-UA-Platform-Version": `"6.6.0"`,
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"chrome-android",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
		"Android",
		true,
		map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.9",
			"Sec-CH-UA":                 `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			"Sec-CH-UA-Mobile":          "?1",
			"Sec-CH-UA-Platform":        `"Android"`,
			"Sec-CH-UA-Platform-Version": `"14.0.0"`,
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"firefox-linux",
		"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Linux",
		false,
		map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"safari-macos",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
		"macOS",
		false,
		map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"safari-iphone",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
		"iOS",
		true,
		map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	))

	UserAgents = append(UserAgents, NewUserAgent(
		"curl",
		"curl/8.7.1",
		"",
		false,
		map[string]string{
			"Accept": "*/*",
		},
	))

}

