package cli

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"

func findAgentBinaryIn(directory string) (string, error) {

	override := os.Getenv("EXOCOMP_AGENT")

	if override != "" {

		stat, err0 := os.Stat(override)

		if err0 == nil && stat.Mode().IsRegular() == true {
			return absoluteAgentPath(override)
		} else {
			return "", fmt.Errorf("cli.FindAgentBinary: exocomp-agent is not installed.")
		}

	} else {

		name := "exocomp-agent"

		if runtime.GOOS == "windows" {
			name = name + ".exe"
		}

		path := filepath.Join(directory, name)
		stat, err1 := os.Stat(path)

		if err1 == nil && stat.Mode().IsRegular() == true {
			return absoluteAgentPath(path)
		} else {

			found, err2 := exec.LookPath("exocomp-agent")

			if err2 == nil {
				return absoluteAgentPath(found)
			} else {
				return "", fmt.Errorf("cli.FindAgentBinary: exocomp-agent is not installed.")
			}

		}

	}

}
