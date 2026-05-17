//go:build windows

package cli

import "errors"
import "os"
import "os/exec"
import "path/filepath"

func FindBrowser() (string, error) {

	candidates := []string{
		"chrome.exe",
		"msedge.exe",
		"chromium.exe",
	}

	for _, candidate := range candidates {

		path, err := exec.LookPath(candidate)

		if err == nil {
			return path, nil
		}

	}

	// 2. Common install locations fallback
	program_files     := os.Getenv("ProgramFiles")
	program_files_x86 := os.Getenv("ProgramFiles(x86)")
	local_app_data    := os.Getenv("LocalAppData")

	paths := []string{
		filepath.Join(program_files,     "Google/Chrome/Application/chrome.exe"),
		filepath.Join(program_files_x86, "Google/Chrome/Application/chrome.exe"),
		filepath.Join(program_files,     "Microsoft/Edge/Application/msedge.exe"),
		filepath.Join(local_app_data,    "Chromium/Application/chrome.exe"),
	}

	for _, path := range paths {

		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

	}

	return "", fmt.Errorf("Error: No supported Browser (Chrome/Edge/Chromium) installed.")

}
