package assets

import "slices"
import "sort"

func mergeStringList(left []string, right []string) []string {

	result := append([]string{}, left...)

	for _, value := range right {

		if value == "" || slices.Contains(result, value) == true {
			continue
		}

		result = append(result, value)

	}

	sort.Strings(result)

	return result

}
