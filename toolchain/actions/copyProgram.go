package actions

import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "path/filepath"

func copyProgram(source string, target string) error {

	err0 := utils.CopyFile(source, target)

	if err0 == nil {

		err1 := os.Chmod(target, 0755)

		if err1 == nil {
			return nil
		} else {
			return fmt.Errorf("actions.copyProgram: Cannot set permissions for %s.", filepath.Base(target))
		}

	} else {
		return fmt.Errorf("actions.copyProgram: Cannot copy %s.", filepath.Base(target))
	}

}
