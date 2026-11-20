package webapp

import "fmt"
import "os"
import "path/filepath"

func CopyWebAssets(base_dir string, public_dir string, design_dir string) error {

	go_root, err0 := GetGoRoot()

	if err0 == nil {

		wasm_exec_source, err1 := findWasmExec(go_root)

		if err1 == nil {

			assets := []struct {
				source       string
				target       string
				is_directory bool
			}{
				{wasm_exec_source, filepath.Join(public_dir, "wasm_exec.js"), false},
				{filepath.Join(base_dir, "app", "public", "index.html"), filepath.Join(public_dir, "index.html"), false},
				{filepath.Join(base_dir, "app", "public", "wasm_init.js"), filepath.Join(public_dir, "wasm_init.js"), false},
				{filepath.Join(base_dir, "app", "public", "app"), filepath.Join(public_dir, "app"), true},
				{design_dir, filepath.Join(public_dir, "design"), true},
			}

			for _, asset := range assets {

				err2 := copyAsset(asset.source, asset.target, asset.is_directory)

				if err2 == nil {

					relative, _ := filepath.Rel(base_dir, asset.target)
					fmt.Fprintf(os.Stdout, "-> Copied %s\n", relative)

				} else {
					return err2
				}

			}

			return nil

		} else {
			return err1
		}

	} else {
		return err0
	}

}
