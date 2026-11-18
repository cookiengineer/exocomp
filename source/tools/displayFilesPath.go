package tools

import "path/filepath"

// displayFilesPath returns a sandbox-relative path when possible, otherwise an
// absolute path for paths granted through capabilities.
func displayFilesPath(sandbox string, path string) string {

	relative, err0 := sanitizeSandboxPath(sandbox, path)

	if err0 == nil {
		return relative
	} else if filepath.IsAbs(path) == true {
		return filepath.Clean(path)
	} else {
		return path
	}

}
