package slices

// UniqueStrings removes duplicates and empty values while preserving the
// original order.
func UniqueStrings(values []string) []string {

	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {

		if value != "" && seen[value] == false {

			seen[value] = true
			result = append(result, value)

		}

	}

	return result

}
