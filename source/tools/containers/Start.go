package containers

import utils_podman "exocomp/utils/podman"
import "fmt"

func (tool *Containers) Start(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Start(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Start: Container \"%s\" started.", name), nil
			} else {
				return "", fmt.Errorf("containers.Start: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Start: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Start: Invalid container name \"%s\".", name)
	}

}
