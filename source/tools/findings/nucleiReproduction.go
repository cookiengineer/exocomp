package findings

import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func nucleiReproduction(line nucleiImportLine, endpoint string) string {

	parts := make([]string, 0)

	if line.TemplateID != "" {
		parts = append(parts, fmt.Sprintf("Nuclei template \"%s\" matched.", line.TemplateID))
	} else {
		parts = append(parts, "Nuclei template matched.")
	}

	if endpoint != "" {
		parts = append(parts, fmt.Sprintf("Endpoint: %s", endpoint))
	}

	if line.MatcherName != "" {
		parts = append(parts, fmt.Sprintf("Matcher: %s", line.MatcherName))
	}

	if len(line.ExtractedResults) > 0 {
		parts = append(parts, fmt.Sprintf("Extracted: %s", strings.Join(line.ExtractedResults, ", ")))
	}

	return utils_findings.LimitText(utils_findings.RedactSecrets(strings.Join(parts, " ")), findingsMaxText)

}
