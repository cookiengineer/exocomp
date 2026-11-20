package containers

import utils_podman "exocomp/utils/podman"
import "fmt"

func (tool *Containers) Create(image string, name string, command string, arguments []string, network string) (string, error) {

	if tool.isAllowedImage(image) == true {

		if tool.isValidName(name) == true {

			if network == "bridge" || network == "none" {

				if command == "" && len(arguments) > 0 {
					return "", fmt.Errorf("containers.Create: Invalid parameter \"arguments\" requires \"command\".")
				}

				if command == "" {
					command = "sleep"
					arguments = []string{"infinity"}
				} else if tool.isValidCommand(command) == false {
					return "", fmt.Errorf("containers.Create: Invalid command \"%s\".", command)
				}

				err1 := utils_podman.Create(image, name, command, arguments, network, tool.Sandbox, tool.label())

				if err1 == nil {
					return fmt.Sprintf("containers.Create: Container \"%s\" created from image \"%s\".", name, image), nil
				} else {
					return "", fmt.Errorf("containers.Create: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("containers.Create: Invalid network \"%s\".", network)
			}

		} else {
			return "", fmt.Errorf("containers.Create: Invalid container name \"%s\".", name)
		}

	} else {
		return "", fmt.Errorf("containers.Create: Invalid image \"%s\": Attempt to use unallowed image.", image)
	}

}
