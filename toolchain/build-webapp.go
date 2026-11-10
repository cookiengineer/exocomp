package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils.GetRoot()

	if err0 == nil {

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build Web App ==\n")

		actions.BuildWebApp(base_dir)

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)
	}

}
