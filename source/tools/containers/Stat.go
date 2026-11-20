package containers

import utils_podman "exocomp/utils/podman"
import "fmt"
import "strings"

func (tool *Containers) Stat(name string) (string, error) {

	if tool.isValidName(name) == true {

		container, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			id := container.ID

			if len(id) > 12 {
				id = id[0:12]
			}

			image := container.ImageName

			if image == "" {
				image = container.Config.Image
			}

			return strings.Join([]string{
				fmt.Sprintf("containers.Stat: %s", strings.TrimPrefix(container.Name, "/")),
				fmt.Sprintf("ID: %s", id),
				fmt.Sprintf("Image: %s", image),
				fmt.Sprintf("Status: %s", container.State.Status),
				fmt.Sprintf("Running: %t", container.State.Running),
				fmt.Sprintf("StartedAt: %s", container.State.StartedAt),
			}, "\n"), nil

		} else {
			return "", fmt.Errorf("containers.Stat: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Stat: Invalid container name \"%s\".", name)
	}

}
