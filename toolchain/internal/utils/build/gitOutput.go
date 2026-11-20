package build

import "os/exec"
import "strings"

func gitOutput(cwd string, arguments ...string) (string, error) {

	cmd := exec.Command("git", arguments...)
	cmd.Dir = cwd

	buffer, err := cmd.Output()

	if err == nil {
		return strings.TrimSpace(string(buffer)), nil
	} else {
		return "", err
	}

}
