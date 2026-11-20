//go:build !windows

package shell

import "fmt"
import "strings"
import "syscall"

// SendSignalByName sends a named signal to the session's process group.
func SendSignalByName(pid int, name string) (string, error) {

	signal := syscall.Signal(0)

	switch strings.ToUpper(strings.TrimSpace(name)) {

	case "INT":
		signal = syscall.SIGINT

	case "TERM":
		signal = syscall.SIGTERM

	case "HUP":
		signal = syscall.SIGHUP

	case "QUIT":
		signal = syscall.SIGQUIT

	case "KILL":
		signal = syscall.SIGKILL

	case "TSTP":
		signal = syscall.SIGTSTP

	case "CONT":
		signal = syscall.SIGCONT

	case "USR1":
		signal = syscall.SIGUSR1

	case "USR2":
		signal = syscall.SIGUSR2

	default:
		return "", fmt.Errorf("shell.SendSignalByName: Invalid signal \"%s\".", name)

	}

	if pid > 0 {

		err := syscall.Kill(-pid, signal)
		if err == nil {
			return strings.ToUpper(strings.TrimSpace(name)), nil
		} else {
			return "", err
		}

	} else {
		return "", fmt.Errorf("shell.SendSignalByName: Invalid process.")
	}

}
