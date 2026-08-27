package yaml

import "strings"

func parseYAMLTag(raw_tag string) string {

	if raw_tag != "" {

		tag_parts := strings.Split(raw_tag, ",")

		return strings.TrimSpace(tag_parts[0])

	} else {
		return ""
	}

}
