package types

import "strings"

func GetUserAgent(name string) *UserAgent {

	search := strings.TrimSpace(strings.ToLower(name))

	for _, useragent := range UserAgents {

		if useragent != nil && strings.ToLower(useragent.Name) == search {
			return useragent
		}

	}

	return nil

}
