package actions

import "fmt"
import "os"
import "path/filepath"

func Clean(base_dir string) error {

	ClearAssets(base_dir)

	build_dir := filepath.Join(base_dir, "build")
	public_dir := filepath.Join(base_dir, "source", "ui", "web", "public")

	operating_systems := []string{"linux", "darwin", "windows"}

	for _, operating_system := range operating_systems {

		target := filepath.Join(build_dir, operating_system)
		err0 := os.RemoveAll(target)

		if err0 == nil {
			path, _ := filepath.Rel(base_dir, target)
			fmt.Fprintf(os.Stdout, "-> Removed %s\n", path)
		} else {
			return err0
		}

	}

	generated := []string{
		filepath.Join(public_dir, "app"),
		filepath.Join(public_dir, "design"),
		filepath.Join(public_dir, "index.html"),
		filepath.Join(public_dir, "main.wasm"),
		filepath.Join(public_dir, "wasm_exec.js"),
		filepath.Join(public_dir, "wasm_init.js"),
	}

	for _, path := range generated {

		err1 := os.RemoveAll(path)

		if err1 != nil {
			return err1
		}

		relative, _ := filepath.Rel(base_dir, path)
		fmt.Fprintf(os.Stdout, "-> Removed %s\n", relative)

	}

	return nil

}
