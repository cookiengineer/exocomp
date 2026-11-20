package websites

import "exocomp/types"
import "fmt"
import "sort"
import "strings"

func (tool *Websites) List() (string, error) {

	names := make([]string, 0)

	for _, useragent := range types.UserAgents {
		names = append(names, useragent.Name)
	}

	sort.Strings(names)

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("websites.List: %d User-Agent presets available.", len(names)))

	for _, name := range names {

		useragent := types.GetUserAgent(name)

		if useragent != nil {
			lines = append(lines, fmt.Sprintf("- Name: \"%s\", Platform: %s, Mobile: %t", useragent.Name, useragent.Platform, useragent.Mobile))
		}

	}

	return strings.Join(lines, "\n"), nil

}
