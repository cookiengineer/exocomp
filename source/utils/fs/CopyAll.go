package fs

import "fmt"
import "os"
import "path/filepath"

func CopyAll(source string, target string) error {

	stat, err := os.Lstat(source)

	if err == nil {

		if stat.Mode()&os.ModeSymlink != 0 {

			return fmt.Errorf("Symbolic links are not allowed.")

		} else if stat.IsDir() {

			err0 := filepath.WalkDir(source, func(path string, entry os.DirEntry, err1 error) error {

				if err1 == nil {

					if entry.Type()&os.ModeSymlink != 0 {
						return nil
					}

					relative, err2 := filepath.Rel(source, path)

					if err2 == nil {

						resolved := filepath.Join(target, relative)

						if entry.IsDir() {
							return os.MkdirAll(resolved, 0755)
						} else {
							return Copy(path, resolved)
						}

					} else {
						return err2
					}

				} else {
					return err1
				}

			})

			return err0

		} else {
			return Copy(source, target)
		}

	} else {
		return err
	}

}
