package main

import "exocomp-toolchain/internal/utils/root"
import "exocomp-toolchain/internal/utils/test"
import "fmt"
import "os"

func main() {

	base_dir, err0 := root.GetRoot()

	if err0 != nil {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)

	}

	err1 := test.Run(base_dir)

	if err1 != nil {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
		os.Exit(1)

	}

}
