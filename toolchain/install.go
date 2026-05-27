package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils.GetRoot()

	if err0 == nil {
		actions.InstallExocomp(base_dir)
	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}
