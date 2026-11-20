package lint

import "strings"

func isToolPackage(relative string) (string, bool) {

	if strings.HasPrefix(relative, "tools/") == false {
		return "", false
	} else {

		parts := strings.Split(relative, "/")

		if len(parts) < 3 {
			return "", false
		} else {
			return parts[1], true
		}

	}

}
