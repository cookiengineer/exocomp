package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils.GetRoot()

	if err0 == nil {

		err1 := actions.Clean(base_dir)

		if err1 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
			os.Exit(1)
		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)
	}

}
