package actions

import "fmt"
import "os"
import "path/filepath"

func ClearAssets(base_dir string) {

	exocomps_dir   := filepath.Join(base_dir, "installer", "assets", "exocomps")
	entries1, err1 := os.ReadDir(exocomps_dir)

	if err1 == nil {

		for _, entry := range entries1 {

			name := entry.Name()

			if name == ".gitkeep" {

				// Do Nothing

			} else {

				binary  := filepath.Join(exocomps_dir, name)
				err     := os.Remove(binary)
				path, _ := filepath.Rel(base_dir, binary)

				if err == nil {
					fmt.Fprintf(os.Stdout, "-> Removed %s\n", path)
				} else {
					fmt.Fprintf(os.Stderr, "!! Error removing %s: %s\n", path, err.Error())
				}

			}

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
	}

	programs_dir   := filepath.Join(base_dir, "installer", "assets", "programs")
	entries2, err2 := os.ReadDir(programs_dir)

	if err2 == nil {

		for _, entry := range entries2 {

			name := entry.Name()

			if name == ".gitkeep" {

				// Do Nothing

			} else {

				binary  := filepath.Join(programs_dir, name)
				err     := os.Remove(binary)
				path, _ := filepath.Rel(base_dir, binary)

				if err == nil {
					fmt.Fprintf(os.Stdout, "-> Removed %s\n", path)
				} else {
					fmt.Fprintf(os.Stderr, "!! Error removing %s: %s\n", path, err.Error())
				}

			}

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err2.Error())
	}

}
