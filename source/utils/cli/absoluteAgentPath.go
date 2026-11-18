package cli

import "fmt"
import "path/filepath"

func absoluteAgentPath(path string) (string, error) {

	absolute, err := filepath.Abs(path)

	if err == nil {
		return absolute, nil
	} else {
		return "", fmt.Errorf("cli.FindAgentBinary: exocomp-agent is not installed.")
	}

}
