package files

import "fmt"
import "os"
import "path/filepath"

func CopyDir(source string, target string) error {

	stat, err1 := os.Stat(source)

	if err1 == nil {

		if stat.IsDir() == true {

			err2 := os.MkdirAll(target, stat.Mode())

			if err2 == nil {

				entries, err3 := os.ReadDir(source)

				if err3 == nil {

					for _, entry := range entries {

						src := filepath.Join(source, entry.Name())
						dst := filepath.Join(target, entry.Name())

						if entry.IsDir() == true {

							err4 := CopyDir(src, dst)

							if err4 != nil {
								return err4
							}

						} else {

							err4 := CopyFile(src, dst)

							if err4 != nil {
								return err4
							}

						}

					}

					return nil

				} else {
					return err3
				}

			} else {
				return err2
			}

		} else {
			return fmt.Errorf("files.CopyDir: %s is not a directory", source)
		}

	} else {
		return err1
	}

}
