package capabilities

import "path/filepath"

// IsCanonicalPath reports whether path is absolute, clean, and already in its
// symlink-resolved canonical form.
func IsCanonicalPath(path string) bool {

	if path == "" || filepath.IsAbs(path) == false {
		return false
	}

	if filepath.Clean(path) != path {
		return false
	}

	return CanonicalPath(path) == path

}
