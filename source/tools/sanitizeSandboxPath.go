package tools

import "fmt"
import "os"
import "path/filepath"

func sanitizeSandboxPath(sandbox string, file_path string) (string, error) {

	if file_path == "" || file_path == "." || file_path == "./" {
		return ".", nil
	} else if file_path == sandbox {
		return ".", nil
	}

	relative, err0 := checkSandboxPath(sandbox, file_path)

	if err0 == nil {

		sandbox_root, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			err2 := checkSandboxSymlinks(sandbox_root, relative)

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
			return "", fmt.Errorf("Invalid path \"%s\": Invalid sandbox \"%s\"", file_path, sandbox)
		}

	} else {
		return "", err0
	}

}
