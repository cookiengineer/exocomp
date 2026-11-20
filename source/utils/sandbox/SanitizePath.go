package sandbox

import "fmt"
import "os"
import "path/filepath"

// SanitizePath validates file_path and returns it relative to sandbox, using a
// leading "./" for display. It does not create directories.
func SanitizePath(sandbox string, file_path string) (string, error) {

	if file_path == "" || file_path == "." || file_path == "./" {
		return ".", nil
	} else if file_path == sandbox {
		return ".", nil
	}

	relative, err0 := CheckPath(sandbox, file_path)

	if err0 == nil {

		sandbox_root, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			err2 := CheckSymlinks(sandbox_root, relative)

			if err2 == nil {

				if relative == "." || relative == "" {
					return "." + string(os.PathSeparator), nil
				} else {
					return "." + string(os.PathSeparator) + relative, nil
				}

			} else {
				return "", err2
			}

		} else {
			return "", fmt.Errorf("sandbox.SanitizePath: Invalid path \"%s\": Invalid sandbox \"%s\"", file_path, sandbox)
		}

	} else {
		return "", err0
	}

}
