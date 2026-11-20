package main

import utils_root "exocomp-toolchain/internal/utils/root"
import "exocomp-toolchain/internal/utils/webapp"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils_root.GetRoot()

	if err0 == nil {

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build Web App ==\n")

		webapp.BuildWebApp(base_dir)

	} else {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)

	}

}
