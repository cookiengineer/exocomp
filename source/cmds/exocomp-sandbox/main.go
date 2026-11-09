package main

import "exocomp/utils/landlock"
import "os"

func main() {

	os.Exit(landlock.Execute(os.Args[1:]))

}
