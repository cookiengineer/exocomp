package tools

import "fmt"

// shellFormatExitCode returns an optional exit-code suffix for tool output.
func shellFormatExitCode(status string, exit_code int) string {

	if status == "running" || status == "starting" || status == "closed" {
		return ""
	} else if exit_code < 0 {
		return ""
	} else {
		return fmt.Sprintf(", exit-code %d", exit_code)
	}

}
