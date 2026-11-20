package containers

import utils_podman "exocomp/utils/podman"
import "fmt"
import "sort"
import "strings"

func (tool *Containers) List() (string, error) {

	containers, err1 := utils_podman.Processes(tool.label())

	if err1 == nil {

		lines := make([]string, 0)

		for _, container := range containers {
			lines = append(lines, fmt.Sprintf("- Name: %s, Image: %s, Status: %s", strings.Join(container.Names, ","), container.Image, container.Status))
		}

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("containers.List: %d containers available.", len(lines)))

		for _, line := range lines {
			result = append(result, line)
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("containers.List: %s", err1.Error())
	}

}
