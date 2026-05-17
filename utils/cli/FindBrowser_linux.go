//go:build linux

package cli

import "fmt"
import "os/exec"

func FindBrowser() (string, error) {

	candidates := []string{
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
	}

	for _, candidate := range candidates {

		path, err := exec.LookPath(candidate)

		if err == nil {
			return path, nil
		}

	}

	return "", fmt.Errorf("Error: No supported Browser (Chromium/Chrome) installed.")

}
