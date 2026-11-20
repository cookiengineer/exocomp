package sandbox

import "fmt"
import "os"
import "path/filepath"

// ResolvePath validates file_path and creates missing parent directories.
// It returns an absolute, symlink-free path inside sandbox.
func ResolvePath(sandbox string, file_path string) (string, error) {

	if file_path == "" || file_path == "." || file_path == "./" {
		return sandbox, nil
	} else if file_path == sandbox {
		return sandbox, nil
	}

	relative, err0 := CheckPath(sandbox, file_path)

	if err0 == nil {

		sandbox_root, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			resolved_path := filepath.Join(sandbox_root, relative)
			parent_folder := filepath.Dir(resolved_path)
			err2 := os.MkdirAll(parent_folder, 0755)

			if err2 == nil {

				err3 := CheckSymlinks(sandbox_root, relative)

				if err3 == nil {
					return resolved_path, nil
				} else {
					return "", err3
				}

			} else {
				return "", fmt.Errorf("sandbox.ResolvePath: Invalid path \"%s\": Cannot create parent folder", relative)
			}

		} else {
			return "", fmt.Errorf("sandbox.ResolvePath: Invalid path \"%s\": Invalid sandbox \"%s\"", file_path, sandbox)
		}

	} else {
		return "", err0
	}

}
