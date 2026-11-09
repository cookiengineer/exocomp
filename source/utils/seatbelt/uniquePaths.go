package seatbelt

func uniquePaths(paths []string) []string {

	seen := make(map[string]bool, len(paths))
	result := make([]string, 0, len(paths))

	for _, path := range paths {

		if path != "" && seen[path] == false {

			seen[path] = true
			result = append(result, path)

		}

	}

	return result

}
