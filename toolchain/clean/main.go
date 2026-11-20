package main

import utils_root "exocomp-toolchain/internal/utils/root"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils_root.GetRoot()

	if err0 == nil {

		err1 := Clean(base_dir)

		if err1 != nil {

			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
			os.Exit(1)

		}

	} else {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)

	}

}
