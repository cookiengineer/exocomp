package actions

import "installer/assets"
import "fmt"
import "io/fs"
import "os"
import "path/filepath"

func InstallPrograms(prefix string) error {

	err0 := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)

	if err0 == nil {

		entries, err1 := fs.ReadDir(assets.FS, "programs")

		if err1 == nil {

			for _, entry := range entries {

				name := entry.Name()
				buffer1, err12 := fs.ReadFile(assets.FS, "programs/" + name)

				if err12 == nil {

					path  := filepath.Join(prefix, "bin", name)
					err13 := os.WriteFile(path, buffer1, 0755)

					if err13 == nil {
						fmt.Fprintf(os.Stdout, "-> Installed %s\n", path)
					} else {
						return err13
					}

				} else {
					return err12
				}

			}

		} else {
			return err1
		}

	} else {
		return err0
	}

	return nil

}
