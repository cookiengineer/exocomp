package containers

import utils_podman "exocomp/utils/podman"
import "fmt"
import "sort"
import "strings"

func (tool *Containers) Images() (string, error) {

	images, err1 := utils_podman.Images()

	if err1 == nil {

		names := make([]string, 0)

		for _, image := range images {

			for _, tag := range image.RepoTags {

				if tool.isAllowedImage(tag) == true {
					names = append(names, tag)
				}

			}

		}

		sort.Strings(names)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("containers.Images: %d images available.", len(names)))

		for _, name := range names {
			result = append(result, fmt.Sprintf("Name: %s", name))
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("containers.Images: %s", err1.Error())
	}

}
