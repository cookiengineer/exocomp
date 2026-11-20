package capabilities

import "path/filepath"

// IsCanonicalPath reports whether path is absolute, clean, and already in its
// symlink-resolved canonical form.
func IsCanonicalPath(path string) bool {

	if path != "" && filepath.IsAbs(path) == true {

		if filepath.Clean(path) == path {
			return CanonicalPath(path) == path
		} else {
			return false
		}
	} else {
		return false
	}
}
