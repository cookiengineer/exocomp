package main

import "exocomp/utils/sandbox"
import "os"

func main() {

	os.Exit(sandbox.Execute(os.Args[1:]))

}
