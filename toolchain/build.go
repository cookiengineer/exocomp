package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils.GetRoot()

	if err0 == nil {

		fmt.Fprint(os.Stdout, "== Clone programs ==\n")

		actions.ClonePrograms(base_dir)

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build linux ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "linux")
		actions.BuildPrograms(base_dir, "linux")
		actions.BuildInstaller(base_dir, "linux")

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build darwin ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "darwin")
		actions.BuildPrograms(base_dir, "darwin")
		actions.BuildInstaller(base_dir, "darwin")

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build windows ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "windows")
		actions.BuildPrograms(base_dir, "windows")
		actions.BuildInstaller(base_dir, "windows")

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}

