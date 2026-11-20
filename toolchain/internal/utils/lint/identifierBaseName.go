package lint

import "strings"

// identifierBaseName is the file name without platform suffixes.
func identifierBaseName(file *File) string {

	base := fileBaseName(file)
	parts := strings.Split(base, "_")

	for len(parts) > 1 {

		last := parts[len(parts)-1]

		if platformParts[last] == true {
			parts = parts[:len(parts)-1]
		} else {
			break
		}

	}

	return strings.Join(parts, "_")

}
