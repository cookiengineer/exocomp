package lint

import "os"
import "path/filepath"
import "strings"

func readModuleName(root string) string {

	source, err0 := os.ReadFile(filepath.Join(root, "go.mod"))

	if err0 == nil {

		for _, line := range strings.Split(string(source), "\n") {

			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, "module ") == true {
				return strings.TrimSpace(strings.TrimPrefix(trimmed, "module "))
			}

		}

		return ""

	} else {
		return ""
	}

}
