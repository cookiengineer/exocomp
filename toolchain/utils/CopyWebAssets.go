package utils

import "fmt"
import "os"
import "path/filepath"

func CopyWebAssets(base_dir string, public_dir string, design_dir string) error {

	go_root, err1 := GetGoRoot()

	if err1 != nil {
		return err1
	}

	wasm_exec_source := filepath.Join(go_root, "lib", "wasm", "wasm_exec.js")
	stat1, err2 := os.Stat(wasm_exec_source)

	if err2 != nil {

		wasm_exec_source = filepath.Join(go_root, "misc", "wasm", "wasm_exec.js")
		stat1, err2 = os.Stat(wasm_exec_source)

	}

	if err2 != nil || stat1.IsDir() == true {
		return fmt.Errorf("CopyWebAssets: wasm_exec.js not found in %s", go_root)
	}

	wasm_exec_target := filepath.Join(public_dir, "wasm_exec.js")
	err3 := CopyFile(wasm_exec_source, wasm_exec_target)

	if err3 != nil {
		return err3
	}

	path1, _ := filepath.Rel(base_dir, wasm_exec_target)
	fmt.Fprintf(os.Stdout, "-> Copied %s\n", path1)

	index_source := filepath.Join(base_dir, "app", "public", "index.html")
	index_target := filepath.Join(public_dir, "index.html")
	err4 := CopyFile(index_source, index_target)

	if err4 != nil {
		return err4
	}

	path2, _ := filepath.Rel(base_dir, index_target)
	fmt.Fprintf(os.Stdout, "-> Copied %s\n", path2)

	init_source := filepath.Join(base_dir, "app", "public", "wasm_init.js")
	init_target := filepath.Join(public_dir, "wasm_init.js")
	err5 := CopyFile(init_source, init_target)

	if err5 != nil {
		return err5
	}

	path3, _ := filepath.Rel(base_dir, init_target)
	fmt.Fprintf(os.Stdout, "-> Copied %s\n", path3)

	app_source := filepath.Join(base_dir, "app", "public", "app")
	app_target := filepath.Join(public_dir, "app")
	err6 := CopyDir(app_source, app_target)

	if err6 != nil {
		return err6
	}

	path4, _ := filepath.Rel(base_dir, app_target)
	fmt.Fprintf(os.Stdout, "-> Copied %s\n", path4)

	design_target := filepath.Join(public_dir, "design")
	err7 := CopyDir(design_dir, design_target)

	if err7 != nil {
		return err7
	}

	path5, _ := filepath.Rel(base_dir, design_target)
	fmt.Fprintf(os.Stdout, "-> Copied %s\n", path5)

	return nil

}
