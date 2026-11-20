package containers

import utils_podman "exocomp/utils/podman"
import "fmt"
import "strings"

func (tool *Containers) Logs(name string, tail int) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			output, err2 := utils_podman.Logs(name, tail)

			if err2 == nil {

				return strings.Join([]string{
					fmt.Sprintf("containers.Logs: %s", name),
					output,
				}, "\n"), nil

			} else {
				return "", fmt.Errorf("containers.Logs: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Logs: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Logs: Invalid container name \"%s\".", name)
	}

}
