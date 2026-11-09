package podman

import "os"
import "os/exec"
import "strings"

func podmanPath() (string, error) {

	override := strings.TrimSpace(os.Getenv("EXOCOMP_PODMAN"))

	if override != "" {
		return override, nil
	}

	path, err := exec.LookPath("podman")

	if err == nil {
		return path, nil
	} else {
		return "", err
	}

}
