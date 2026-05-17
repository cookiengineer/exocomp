//go:build darwin

package cli

import "fmt"
import "os"

func CheckRoot() error {

	if os.Geteuid() == 0 {

		return nil

	} else {

		if sudo_user := os.Getenv("SUDO_USER"); sudo_user != "" {
			exe, _ := os.Executable()
			return fmt.Errorf("Error: Run this program with sudo %s", exe)
		} else {
			return fmt.Errorf("Error: Run this program as root")
		}

	}

}
