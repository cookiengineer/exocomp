package containers

import utils_podman "exocomp/utils/podman"
import "fmt"

func (tool *Containers) Stop(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Stop(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Stop: Container \"%s\" stopped.", name), nil
			} else {
				return "", fmt.Errorf("containers.Stop: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Stop: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Stop: Invalid container name \"%s\".", name)
	}

}
