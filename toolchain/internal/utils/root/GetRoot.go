package root

import "fmt"
import "path/filepath"
import "runtime"

func GetRoot() (string, error) {

	_, self, _, ok := runtime.Caller(0)

	if ok == true {

		// self = ./toolchain/internal/utils/root/GetRoot.go
		root_dir := filepath.Dir(self)
		utils_dir := filepath.Dir(root_dir)
		internal_dir := filepath.Dir(utils_dir)
		toolchain_dir := filepath.Dir(internal_dir)
		root := filepath.Dir(toolchain_dir)

		return root, nil

	} else {
		return "", fmt.Errorf("root.GetRoot: Cannot determine current file.")
	}

}
