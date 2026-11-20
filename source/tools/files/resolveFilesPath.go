package files

import utils_capabilities "exocomp/utils/capabilities"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "path/filepath"

// resolveFilesPath resolves a path inside the sandbox or inside a granted
// capability scope. The internal .exocomp folder is reserved for tool-managed
// state and is never accessible through files.*.
func resolveFilesPath(sandbox string, capabilities *utils_capabilities.Set, path string, is_write bool) (string, error) {

	if utils_sandbox.IsReservedPath(sandbox, path) == false {

		resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

		if err0 == nil {

			return resolved, nil

		} else if capabilities != nil && filepath.IsAbs(path) == true {

			absolute, err1 := filepath.Abs(path)

			if err1 == nil {

				canonical, err2 := filepath.EvalSymlinks(absolute)

				if err2 != nil {

					parent, err3 := filepath.EvalSymlinks(filepath.Dir(absolute))

					if err3 == nil {
						canonical = filepath.Join(parent, filepath.Base(absolute))
					}

				}

				if canonical != "" {

					if utils_sandbox.IsReservedPath(sandbox, canonical) == true {
						return "", fmt.Errorf("Path is reserved for internal state.")
					}

					if is_write == true && capabilities.CanWrite(canonical) == true {
						return canonical, nil
					} else if is_write == false && capabilities.CanRead(canonical) == true {
						return canonical, nil
					}

				}

			}

		}

		return "", err0

	} else {
		return "", fmt.Errorf("Path is reserved for internal state.")
	}
}
