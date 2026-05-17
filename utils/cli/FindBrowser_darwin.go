//go:build darwin

package cli

import "fmt"
import "os"
import "os/exec"
import "path/filepath"

func FindBrowser() (string, error) {

	candidates := []string{
		"chromium",
		"google-chrome",
		"chrome",
		"msedge",
		"microsoft-edge",
	}

	for _, candidate := range candidates {

		path, err := exec.LookPath(candidate)

		if err == nil {
			return path, nil
		}

	}

	app_candidates := []string{
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		filepath.Join(os.Getenv("HOME"), "Applications/Chromium.app/Contents/MacOS/Chromium"),
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		filepath.Join(os.Getenv("HOME"), "Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		filepath.Join(os.Getenv("HOME"), "Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"),
	}

	for _, path := range app_candidates {

		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

	}

	return "", fmt.Errorf("Error: No supported Browser (Chromium/Chrome/Edge) installed.")

}
