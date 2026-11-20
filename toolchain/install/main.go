package main

import utils_root "exocomp-toolchain/internal/utils/root"
import "fmt"
import "os"

func main() {

	base_dir, err0 := utils_root.GetRoot()

	if err0 == nil {
		InstallExocomp(base_dir)
	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}
