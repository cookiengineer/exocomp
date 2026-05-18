package actions

import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "path/filepath"

func BuildInstaller(base_dir string, operating_system string) {

	build_dir     := filepath.Join(base_dir, "build")
	installer_dir := filepath.Join(base_dir, "installer")

	err0 := os.MkdirAll(filepath.Join(build_dir, operating_system), 0755)

	if err0 == nil {

		output  := filepath.Join(build_dir, operating_system, "exocomp-installer")
		err1    := utils.Build(installer_dir, "./cmds/exocomp-installer/main.go", output, []string{}, operating_system)
		path, _ := filepath.Rel(base_dir, output)

		if err1 == nil {
			fmt.Fprintf(os.Stdout, "-> Built %s\n", path)
		} else {
			fmt.Fprintf(os.Stderr, "!! Error building %s: %s\n", path, err1.Error())
		}

	} else {

		if err0 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		}

	}

}
