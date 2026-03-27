package tools

import "fmt"
import "os"
import "path/filepath"

func resolveSandboxPath(sandbox string, file_path string) (string, error) {

	tmp1 := filepath.Join(sandbox, file_path)
	resolved_path, err0 := filepath.Abs(tmp1)

	if err0 == nil {

		sandbox_root, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			relative, err2 := filepath.Rel(sandbox_root, resolved_path)

			if err2 == nil {

				if relative == ".." {
					return "", fmt.Errorf("Invalid file path \"%s\": Path tried to escape sandbox", relative)
				} else if len(relative) >= 3 && relative[0:3] == ".." + string(os.PathSeparator) {
					return "", fmt.Errorf("Invalid file path \"%s\": Path tried to escape sandbox", relative)
				} else {

					parent_folder := filepath.Dir(resolved_path)

					err3 := os.MkdirAll(parent_folder, 0755)

					if err3 == nil {
						return resolved_path, nil
					} else {
						return "", fmt.Errorf("Cannot create parent folder of file path \"%s\": %s", resolved_path, err3.Error())
					}

				}

			} else {
				return "", fmt.Errorf("Cannot resolve relative path \"%s\": %s", resolved_path, err2.Error())
			}

		} else {
			return "", fmt.Errorf("Cannot resolve path of sandbox \"%s\": %s", sandbox, err1.Error())
		}

	} else {
		return "", fmt.Errorf("Cannot resolve file path of \"%s\": %s", tmp1, err0.Error())
	}

}
