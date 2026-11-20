package build

import utils_files "exocomp-toolchain/internal/utils/files"
import "fmt"
import "os"
import "path/filepath"

func copyProgram(source string, target string) error {

	err0 := utils_files.CopyFile(source, target)

	if err0 == nil {

		err1 := os.Chmod(target, 0755)

		if err1 == nil {
			return nil
		} else {
			return fmt.Errorf("build.CopyProgram: Cannot set permissions for %s.", filepath.Base(target))
		}

	} else {
		return fmt.Errorf("build.CopyProgram: Cannot copy %s.", filepath.Base(target))
	}

}
