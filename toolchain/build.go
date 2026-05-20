package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "strings"

func main() {

	base_dir, err0 := utils.GetRoot()

	build_linux   := true
	build_darwin  := true
	build_windows := true

	if len(os.Args) > 1 {

		tmp1 := strings.TrimSpace(os.Args[1])

		switch tmp1 {
		case "linux":
			build_darwin  = false
			build_windows = false
		case "darwin":
			build_linux   = false
			build_windows = false
		case "windows":
			build_linux  = false
			build_darwin = false
		}

	}

	if err0 == nil {

		fmt.Fprint(os.Stdout, "== Clone programs ==\n")

		actions.ClonePrograms(base_dir)

		if build_linux == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build linux ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "linux")
			actions.BuildPrograms(base_dir, "linux")
			actions.BuildInstaller(base_dir, "linux")

		}

		if build_darwin == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build darwin ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "darwin")
			actions.BuildPrograms(base_dir, "darwin")
			actions.BuildInstaller(base_dir, "darwin")

		}

		if build_windows == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build windows ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "windows")
			actions.BuildPrograms(base_dir, "windows")
			actions.BuildInstaller(base_dir, "windows")

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}

