package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Search(arguments utils_chat.Arguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == true {

		query = strings.ToLower(strings.TrimSpace(query))

		if query == "" {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is empty.")
		}

		severity_filter, has_severity := arguments.String("severity")
		status_filter, has_status := arguments.String("status")
		limit := findingsDefaultLimit

		if value, ok := arguments.Int("limit"); ok == true {

			if value > 0 {
				limit = value
			}

		}

		if limit > findingsMaxLimit {
			limit = findingsMaxLimit
		}

		if has_severity == true && types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == false {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"severity\".")
		}

		if has_status == true && types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == false {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"status\".")
		}

		tool.mutex.Lock()
		findings, err0 := utils_findings.ReadAll(tool.Playground)

		if err0 == nil {

			summaries := make([]types.FindingSummary, 0)

			for _, finding := range findings {

				if has_severity == true && string(finding.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
					continue
				} else if has_status == true && string(finding.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
					continue
				} else if utils_findings.MatchesQuery(finding, query) == false {
					continue
				}

				summaries = append(summaries, types.FindingSummaryFromFinding(finding))

			}

			utils_findings.SortSummaries(summaries)

			total := len(summaries)

			if len(summaries) > limit {
				summaries = summaries[0:limit]
			}

			lines := make([]string, 0, len(summaries))

			for _, summary := range summaries {
				lines = append(lines, formatFindingSummary(summary))
			}

			tool.mutex.Unlock()

			result := make([]string, 0)
			result = append(result, fmt.Sprintf("findings.Search: %d matches (showing %d).", total, len(lines)))

			for _, line := range lines {
				result = append(result, line)
			}

			return strings.Join(result, "\n"), nil

		} else {
			tool.mutex.Unlock()
			return "", fmt.Errorf("findings.Search: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is not a string.")
	}

}
