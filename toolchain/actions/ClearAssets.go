package actions

import "fmt"
import "os"
import "path/filepath"

func ClearAssets(base_dir string) {

	assets_dir := filepath.Join(base_dir, "installer", "assets", "exocomps")

	entries, err0 := os.ReadDir(assets_dir)

	if err0 == nil {

		for _, entry := range entries {

			name := entry.Name()

			if name == ".gitkeep" {

				// Do Nothing

			} else {

				binary  := filepath.Join(assets_dir, name)
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

		if err0 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		}

	}

}
