package fs

import "os"
import "path/filepath"

func CopyAll(source string, target string) error {

	stat, err := os.Stat(source)

	if err == nil {

		if stat.IsDir() {

			err0 := filepath.WalkDir(source, func(path string, entry os.DirEntry, err1 error) error {

				if err1 == nil {

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
