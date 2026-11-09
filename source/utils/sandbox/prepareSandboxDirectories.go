package sandbox

import "errors"
import "fmt"
import "io/fs"
import "os"

func prepareSandboxDirectories(directories []string) error {

	for _, directory := range directories {

		err := os.MkdirAll(directory, 0755)

		if err != nil {

			if errors.Is(err, fs.ErrPermission) {
				return fmt.Errorf("sandbox.prepareSandboxDirectories: Permission denied.")
			} else {
				return fmt.Errorf("sandbox.prepareSandboxDirectories: Cannot prepare sandbox.")
			}

		}

	}

	return nil

}
