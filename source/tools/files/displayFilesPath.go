package files

import utils_sandbox "exocomp/utils/sandbox"
import "path/filepath"

// displayFilesPath returns a sandbox-relative path when possible, otherwise an
// absolute path for paths granted through capabilities.
func displayFilesPath(sandbox string, path string) string {

	relative, err0 := utils_sandbox.SanitizePath(sandbox, path)

	if err0 == nil {
		return relative
	} else if filepath.IsAbs(path) == true {
		return filepath.Clean(path)
	} else {
		return path
	}

}
