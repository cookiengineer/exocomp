package main

import "exocomp-toolchain/actions"
import "fmt"
import "os"
import "path/filepath"
import "runtime"

func getRepositoryRoot() (string, error) {

	_, self, _, ok := runtime.Caller(0)

	if ok == true {

		// self = ./toolchain/build.go
		toolchain := filepath.Dir(self)

		// root = parent of toolchain/
		root := filepath.Dir(toolchain)

		return root, nil

	} else {
		return "", fmt.Errorf("unable to determine current file")
	}

}

func main() {

	base_dir, err0 := getRepositoryRoot()

	if err0 == nil {

		fmt.Fprint(os.Stdout, "== Build linux ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "linux")
		actions.BuildInstaller(base_dir, "linux")

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build darwin ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "darwin")
		actions.BuildInstaller(base_dir, "darwin")

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build windows ==\n")

		actions.ClearAssets(base_dir)
		actions.BuildExocomps(base_dir, "windows")
		actions.BuildInstaller(base_dir, "windows")

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}

