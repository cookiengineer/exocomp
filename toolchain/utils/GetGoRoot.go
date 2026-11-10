package utils

import "fmt"
import "os/exec"
import "strings"

func GetGoRoot() (string, error) {

	cmd := exec.Command("go", "env", "GOROOT")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	goroot := strings.TrimSpace(string(output))

	if goroot == "" {
		return "", fmt.Errorf("GetGoRoot: GOROOT is empty")
	}

	return goroot, nil

}
