package capabilities

import "os"
import "strings"

// IsPathDenied reports whether a path is on the capability denylist.
func IsPathDenied(path string) bool {

	for _, denied := range deniedPaths {

		if path == denied {
			return true
		} else if denied != "/" && strings.HasPrefix(path, denied+string(os.PathSeparator)) == true {
			return true
		}

	}

	return false

}
