package tools

import "fmt"
import "os"
import "path/filepath"
import "strings"

func checkSandboxSymlinks(sandbox_root string, relative string) error {

	if relative == "." || relative == "" {
		return nil
	}

	current := sandbox_root

	for _, part := range strings.Split(relative, string(os.PathSeparator)) {

		current = filepath.Join(current, part)

		info, err := os.Lstat(current)

		if err == nil {

			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("Invalid path \"%s\": Symbolic links are not allowed", "."+string(os.PathSeparator)+relative)
			}

		} else if os.IsNotExist(err) {
			break
		} else if os.IsPermission(err) {
			return fmt.Errorf("Invalid path \"%s\": Permission denied.", "."+string(os.PathSeparator)+relative)
		} else {
			return fmt.Errorf("Invalid path \"%s\": Cannot access path", "."+string(os.PathSeparator)+relative)
		}

	}

	return nil

}
