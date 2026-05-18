package utils

import "fmt"
import "path/filepath"
import "runtime"

func GetRoot() (string, error) {

	_, self, _, ok := runtime.Caller(0)

	if ok == true {

		// self = ./toolchain/utils/GetRoot.go
		utils := filepath.Dir(self)

		// toolchain = parent of self
		toolchain := filepath.Dir(utils)

		// root = parent of toolchain/
		root := filepath.Dir(toolchain)

		return root, nil

	} else {
		return "", fmt.Errorf("unable to determine current file")
	}

}

