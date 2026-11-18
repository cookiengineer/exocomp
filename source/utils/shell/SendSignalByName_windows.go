//go:build windows

package shell

import "fmt"

// SendSignalByName is not supported on Windows.
func SendSignalByName(pid int, name string) (string, error) {
	return "", fmt.Errorf("Shell signals are not supported on this system.")
}
