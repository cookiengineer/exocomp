package capabilities

import "path/filepath"
import "strings"

// IsPathWithin reports whether path is equal to or below scope.
func IsPathWithin(scope string, path string) bool {

	scope = filepath.Clean(scope)
	path = filepath.Clean(path)

	if path == scope {
		return true
	} else {

		relative, err0 := filepath.Rel(scope, path)

		if err0 == nil {
			if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) == true {
				return false
			} else {
				return true
			}
		}

		return false

	}

}
