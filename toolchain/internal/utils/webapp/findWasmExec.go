package webapp

import "fmt"
import "os"
import "path/filepath"

func findWasmExec(go_root string) (string, error) {

	lib_path := filepath.Join(go_root, "lib", "wasm", "wasm_exec.js")
	misc_path := filepath.Join(go_root, "misc", "wasm", "wasm_exec.js")

	stat1, err1 := os.Stat(lib_path)

	if err1 == nil {

		if stat1.IsDir() == false {
			return lib_path, nil
		} else {
			return "", fmt.Errorf("webapp.CopyWebAssets: wasm_exec.js not found in %s", go_root)
		}

	} else {

		stat2, err2 := os.Stat(misc_path)

		if err2 == nil && stat2.IsDir() == false {
			return misc_path, nil
		} else {
			return "", fmt.Errorf("webapp.CopyWebAssets: wasm_exec.js not found in %s", go_root)
		}

	}

}
