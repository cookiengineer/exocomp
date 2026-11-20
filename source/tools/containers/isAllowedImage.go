package containers

import "slices"
import "strings"

func (tool *Containers) isAllowedImage(image string) bool {

	if image == "" || len(image) > 255 {
		return false
	} else if strings.HasPrefix(image, "-") == true {
		return false
	}

	for _, character := range image {

		if character >= 'a' && character <= 'z' {
			continue
		} else if character >= 'A' && character <= 'Z' {
			continue
		} else if character >= '0' && character <= '9' {
			continue
		} else if strings.ContainsRune("._/:@-", character) == true {
			continue
		} else {
			return false
		}

	}

	return slices.Contains(tool.AllowedImages, image)

}
