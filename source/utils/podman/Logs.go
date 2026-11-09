package podman

import "strconv"
import "time"

func Logs(name string, tail int) (string, error) {

	arguments := []string{
		"logs",
		"--tail", strconv.Itoa(tail),
		name,
	}

	return execute(arguments, 30*time.Second, 0, name, "")

}
