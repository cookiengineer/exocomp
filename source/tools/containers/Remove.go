package containers

import utils_podman "exocomp/utils/podman"
import "fmt"

func (tool *Containers) Remove(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Remove(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Remove: Container \"%s\" removed.", name), nil
			} else {
				return "", fmt.Errorf("containers.Remove: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Remove: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Remove: Invalid container name \"%s\".", name)
	}

}
