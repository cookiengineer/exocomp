package sandbox

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"

// HelperPath returns the absolute path to the exocomp-sandbox helper.
func HelperPath() (string, error) {

	override := os.Getenv("EXOCOMP_SANDBOX")

	if override != "" {

		stat, err0 := os.Stat(override)

		if err0 == nil && stat.Mode().IsRegular() == true {

			absolute, err1 := absolutePath(override)

			if err1 == nil {
				return absolute, nil
			} else {
				return "", err1
			}

		} else {
			return "", fmt.Errorf("sandbox.HelperPath: Sandbox helper is not installed.")
		}

	} else {

		exe, err2 := os.Executable()

		if err2 == nil {

			name := "exocomp-sandbox"

			if runtime.GOOS == "windows" {
				name = name + ".exe"
			}

			path := filepath.Join(filepath.Dir(exe), name)
			stat, err3 := os.Stat(path)

			if err3 == nil && stat.Mode().IsRegular() == true {

				absolute, err4 := absolutePath(path)

				if err4 == nil {
					return absolute, nil
				} else {
					return "", err4
				}

			}

		}

		path, err5 := exec.LookPath("exocomp-sandbox")

		if err5 == nil {

			absolute, err6 := absolutePath(path)

			if err6 == nil {
				return absolute, nil
			} else {
				return "", err6
			}

		} else {
			return "", fmt.Errorf("sandbox.HelperPath: Sandbox helper is not installed.")
		}

	}

}
