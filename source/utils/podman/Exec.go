package podman

import "time"

func Exec(name string, command string, arguments []string, workdir string) (string, error) {

	podman_arguments := []string{
		"exec",
		"--workdir", workdir,
		name,
		command,
	}

	podman_arguments = append(podman_arguments, arguments...)

	return execute(podman_arguments, 10*time.Minute, 1*time.Minute, name, "")

}
