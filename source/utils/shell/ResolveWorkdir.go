package shell

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"
import "path/filepath"

// ResolveWorkdir validates a sandbox-relative working directory.
func ResolveWorkdir(sandbox string, workdir string) (string, error) {

	if workdir == "" {
		workdir = "."
	}

	relative, err0 := utils_sandbox.CheckPath(sandbox, workdir)

	if err0 == nil {

		sandbox_root, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			resolved := filepath.Join(sandbox_root, relative)
			stat, err2 := os.Stat(resolved)

			if err2 == nil && stat.IsDir() == true {
				return resolved, nil
			} else {
				return "", fmt.Errorf("Invalid workdir \"%s\": Directory does not exist.", workdir)
			}

		} else {
			return "", fmt.Errorf("Invalid workdir \"%s\": Invalid sandbox.", workdir)
		}

	} else {
		return "", err0
	}

}
