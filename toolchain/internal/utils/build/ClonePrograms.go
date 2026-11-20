package build

import "fmt"
import "os"
import "path/filepath"

func ClonePrograms(base_dir string) error {

	vendor_dir := filepath.Join(base_dir, "vendor")
	err0 := os.MkdirAll(vendor_dir, 0755)

	if err0 == nil {

		programs, err1 := ReadPrograms(base_dir)

		if err1 == nil {

			seen := make(map[string]bool, len(programs))

			for _, program := range programs {

				if seen[program.Folder] == false {

					seen[program.Folder] = true
					folder := filepath.Join(vendor_dir, program.Folder)
					stat, err2 := os.Stat(folder)

					if err2 == nil && stat.IsDir() == true {
						err2 = CheckoutRepository(folder, program.Version)
					} else {
						err2 = CloneRepository(vendor_dir, program.Repository, folder, program.Version)
					}

					path, _ := filepath.Rel(base_dir, folder)

					if err2 == nil {
						fmt.Fprintf(os.Stdout, "-> Checked out %s at %s\n", path, program.Version)
					} else {
						return fmt.Errorf("build.ClonePrograms: Cannot prepare %s at %s.", path, program.Version)
					}

				}

			}

			return nil

		} else {
			return err1
		}

	} else {
		return fmt.Errorf("build.ClonePrograms: Cannot create vendor folder.")
	}

}
