package actions

import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "path/filepath"
import "runtime"

func InstallExocomp(base_dir string) {

	build_dir := filepath.Join(base_dir, "build")
	programs  := []string{"exocomp", "exocomp-agent", "agimus"}

	if runtime.GOOS == "darwin" {

		for _, program := range programs {

			prefix := "/usr/local/bin"
			build  := filepath.Join(build_dir, "darwin", program)
			target := filepath.Join(prefix, program)

			err := utils.CopyFile(build, target)

			if err != nil {
				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
			}

		}

	} else if runtime.GOOS == "linux" {

		for _, program := range programs {

			prefix := "/usr/bin"
			build  := filepath.Join(build_dir, "linux", program)
			target := filepath.Join(prefix, program)

			err := utils.CopyFile(build, target)

			if err == nil {
				fmt.Fprintf(os.Stdout, "-> Installed %s\n", target)
			} else {
				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
			}

		}

	} else if runtime.GOOS == "windows" {

		for _, program := range programs {

			prefix := "C:\\Program Files"

			if runtime.GOARCH == "386" {
				prefix = "C:\\Program Files (x86)"
			}

			build  := filepath.Join(build_dir, "linux", fmt.Sprintf("%s.exe", program))
			target := filepath.Join(prefix, program)

			err := utils.CopyFile(build, target)

			if err == nil {
				fmt.Fprintf(os.Stdout, "-> Installed %s\n", target)
			} else {
				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
			}

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: Unsupported operating system \"%s\"\n", runtime.GOOS)
	}

}
