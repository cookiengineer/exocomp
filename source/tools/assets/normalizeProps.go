package assets

import utils_findings "exocomp/utils/findings"
import "strings"

func normalizeProps(props map[string]string) map[string]string {

	result := make(map[string]string, len(props))

	for key, value := range props {

		key = strings.TrimSpace(utils_findings.RedactSecrets(key))
		value = strings.TrimSpace(utils_findings.RedactSecrets(value))

		if key == "" || value == "" {
			continue
		}

		result[utils_findings.LimitText(key, assetsMaxText)] = utils_findings.LimitText(value, assetsMaxText)

	}

	return result

}
