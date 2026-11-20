package assets

import utils_findings "exocomp/utils/findings"
import "slices"
import "sort"
import "strings"

func normalizeStringList(values []string, maximum int) []string {

	result := make([]string, 0, len(values))

	for _, value := range values {

		value = strings.TrimSpace(utils_findings.RedactSecrets(value))

		if value == "" || slices.Contains(result, value) == true {
			continue
		} else {
			result = append(result, utils_findings.LimitText(value, maximum))
		}
	}

	sort.Strings(result)

	return result

}
