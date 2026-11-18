//go:build windows

package tools

import "fmt"

// shellSendSignalByName is not supported on Windows.
func shellSendSignalByName(pid int, name string) (string, error) {
	return "", fmt.Errorf("Shell signals are not supported on this system.")
}
