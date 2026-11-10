package actions

import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "path/filepath"

func BuildWebApp(base_dir string) {

	app_dir := filepath.Join(base_dir, "app")
	public_dir := filepath.Join(base_dir, "source", "ui", "web", "public")
	gooey_dir := filepath.Join(filepath.Dir(base_dir), "gooey")
	design_dir := filepath.Join(gooey_dir, "design")

	stat1, err1 := os.Stat(app_dir)

	if err1 == nil && stat1.IsDir() == true {

		stat2, err2 := os.Stat(design_dir)

		if err2 == nil && stat2.IsDir() == true {

			err3 := os.MkdirAll(public_dir, 0755)

			if err3 == nil {

				generated := []string{
					filepath.Join(public_dir, "app"),
					filepath.Join(public_dir, "design"),
					filepath.Join(public_dir, "index.html"),
					filepath.Join(public_dir, "main.wasm"),
					filepath.Join(public_dir, "wasm_exec.js"),
					filepath.Join(public_dir, "wasm_init.js"),
				}

				removed := true

				for _, path := range generated {

					err4 := os.RemoveAll(path)

					if err4 != nil {
						fmt.Fprintf(os.Stderr, "!! Error removing %s: %s\n", path, err4.Error())
						removed = false
						break
					}

				}

				if removed == true {

					wasm_output := filepath.Join(public_dir, "main.wasm")
					err5 := utils.BuildWASM(app_dir, "main.go", wasm_output)

					if err5 == nil {

						path1, _ := filepath.Rel(base_dir, wasm_output)
						fmt.Fprintf(os.Stdout, "-> Built %s\n", path1)

						err6 := utils.CopyWebAssets(base_dir, public_dir, design_dir)

						if err6 == nil {
							// Do Nothing
						} else {
							fmt.Fprintf(os.Stderr, "!! Error copying Web App assets: %s\n", err6.Error())
						}

					} else {
						fmt.Fprintf(os.Stderr, "!! Error building %s: %s\n", wasm_output, err5.Error())
					}

				}

			} else {
				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err3.Error())
			}

		} else {
			fmt.Fprintf(os.Stderr, "!! Error: Gooey design directory %s does not exist.\n", design_dir)
		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: App directory %s does not exist.\n", app_dir)
	}

}
