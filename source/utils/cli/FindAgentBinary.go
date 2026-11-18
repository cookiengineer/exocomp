package cli

import "fmt"
import "os"
import "path/filepath"

// FindAgentBinary returns the absolute path to the exocomp-agent child process
// binary. It checks the EXOCOMP_AGENT override, the folder of the running
// binary, and finally PATH.
func FindAgentBinary() (string, error) {

	exe, err0 := os.Executable()

	if err0 == nil {
		return findAgentBinaryIn(filepath.Dir(exe))
	} else {
		return "", fmt.Errorf("cli.FindAgentBinary: exocomp-agent is not installed.")
	}

}
