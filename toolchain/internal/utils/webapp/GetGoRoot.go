package webapp

import "fmt"
import "os/exec"
import "strings"

func GetGoRoot() (string, error) {

	cmd := exec.Command("go", "env", "GOROOT")
	output, err := cmd.Output()

	if err == nil {

		goroot := strings.TrimSpace(string(output))

		if goroot != "" {
			return goroot, nil
		} else {
			return "", fmt.Errorf("webapp.GetGoRoot: GOROOT is empty.")
		}

	} else {
		return "", err
	}

}
