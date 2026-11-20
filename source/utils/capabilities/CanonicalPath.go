package capabilities

import "path/filepath"

// CanonicalPath returns an absolute cleaned path with symlinks resolved when possible.
func CanonicalPath(path string) string {

	absolute, err0 := filepath.Abs(path)

	if err0 != nil {
		return filepath.Clean(path)
	} else {

		resolved, err1 := filepath.EvalSymlinks(absolute)

		if err1 == nil {
			return resolved
		} else {

			parent := filepath.Dir(absolute)
			resolved_parent, err2 := filepath.EvalSymlinks(parent)

			if err2 == nil {
				return filepath.Join(resolved_parent, filepath.Base(absolute))
			} else {
				return absolute
			}

		}

	}

}
