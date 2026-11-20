package sandbox

import "fmt"
import "os"
import "path/filepath"
import "strings"

// CheckSymlinks rejects symbolic links in the relative path below sandbox_root.
func CheckSymlinks(sandbox_root string, relative string) error {

	if relative == "." || relative == "" {
		return nil
	} else {

		current := sandbox_root

		for _, part := range strings.Split(relative, string(os.PathSeparator)) {

			current = filepath.Join(current, part)

			info, err := os.Lstat(current)

			if err == nil {

				if info.Mode()&os.ModeSymlink != 0 {
					return fmt.Errorf("sandbox.CheckSymlinks: Invalid path \"%s\": Symbolic links are not allowed", "."+string(os.PathSeparator)+relative)
				}

			} else if os.IsNotExist(err) == true {
				break
			} else if os.IsPermission(err) == true {
				return fmt.Errorf("sandbox.CheckSymlinks: Invalid path \"%s\": Permission denied.", "."+string(os.PathSeparator)+relative)
			} else {
				return fmt.Errorf("sandbox.CheckSymlinks: Invalid path \"%s\": Cannot access path", "."+string(os.PathSeparator)+relative)
			}

		}

		return nil

	}
}
