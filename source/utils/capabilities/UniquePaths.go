package capabilities

import "sort"

// UniquePaths removes duplicates and returns the paths sorted.
func UniquePaths(paths []string) []string {

	seen := make(map[string]bool, len(paths))
	result := make([]string, 0, len(paths))

	for _, path := range paths {

		if seen[path] == false {

			seen[path] = true
			result = append(result, path)

		}

	}

	sort.Strings(result)

	return result

}
