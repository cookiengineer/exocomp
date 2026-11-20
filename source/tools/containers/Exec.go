package containers

import utils_podman "exocomp/utils/podman"
import "fmt"
import "strings"

func (tool *Containers) Exec(name string, command string, arguments []string, workdir string) (string, error) {

	if tool.isValidName(name) == true {

		if tool.isValidCommand(command) == true {

			if strings.HasPrefix(workdir, "/") == true && tool.isValidCommand(workdir) == true {

				_, err1 := utils_podman.Inspect(name, tool.label())

				if err1 == nil {

					output, err2 := utils_podman.Exec(name, command, arguments, workdir)

					if err2 == nil {

						command_line := strings.Join(append([]string{command}, arguments...), " ")

						return strings.Join([]string{
							fmt.Sprintf("containers.Exec: podman exec %s %s", name, command_line),
							output,
						}, "\n"), nil

					} else {
						return "", fmt.Errorf("containers.Exec: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("containers.Exec: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("containers.Exec: Invalid workdir \"%s\".", workdir)
			}

		} else {
			return "", fmt.Errorf("containers.Exec: Invalid command \"%s\".", command)
		}

	} else {
		return "", fmt.Errorf("containers.Exec: Invalid container name \"%s\".", name)
	}

}
