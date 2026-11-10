package path

import "strings"

func Collapse(package_path string) string {

	segments := make([]string, 0)

	for _, segment := range strings.Split(package_path, "/") {

		if segment != "" {
			segments = append(segments, segment)
		}

	}

	if len(segments) > 0 {
		return segments[0]
	}

	return ""

}
