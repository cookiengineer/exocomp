package lint

import "strings"

func classifyImport(path string) int {

	if strings.HasPrefix(path, "exocomp/") == true || strings.HasPrefix(path, "exocomp-toolchain/") == true {
		return groupProject
	} else {

		first := path

		if index := strings.Index(path, "/"); index >= 0 {
			first = path[:index]
		}

		if strings.Contains(first, ".") == true {
			return groupThird
		} else {
			return groupStdlib
		}

	}

}
