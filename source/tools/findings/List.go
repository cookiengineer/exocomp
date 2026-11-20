package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "slices"
import "strings"

func (tool *Findings) List(arguments utils_chat.Arguments) (string, error) {

	severity_filter, has_severity := arguments.String("severity")
	status_filter, has_status := arguments.String("status")
	confidence_filter, has_confidence := arguments.String("confidence")
	endpoint_filter, has_endpoint := arguments.String("endpoint")
	asset_filter, has_asset := arguments.String("asset_id")
	limit := findingsDefaultLimit

	if value, ok := arguments.Int("limit"); ok == true {

		if value > 0 {
			limit = value
		}

	}

	if limit > findingsMaxLimit {
		limit = findingsMaxLimit
	}

	if has_severity == false || types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == true {

		if has_status == false || types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == true {

			if has_confidence == false || types.FindingConfidence(strings.ToLower(strings.TrimSpace(confidence_filter))).IsValid() == true {

				tool.mutex.Lock()
				err0 := readFindingsIndex(tool)

				if err0 == nil {

					summaries := make([]types.FindingSummary, 0)

					for _, summary := range tool.contents {

						if has_severity == true && string(summary.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
							continue
						} else if has_status == true && string(summary.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
							continue
						} else if has_confidence == true && string(summary.Confidence) != strings.ToLower(strings.TrimSpace(confidence_filter)) {
							continue
						} else if has_endpoint == true && strings.Contains(strings.ToLower(summary.Endpoint), strings.ToLower(strings.TrimSpace(endpoint_filter))) == false {
							continue
						} else if has_asset == true && slices.Contains(summary.AssetIDs, strings.TrimSpace(asset_filter)) == false {
							continue
						}

						summaries = append(summaries, summary)

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
					result = append(result, fmt.Sprintf("findings.List: %d findings (showing %d).", total, len(lines)))

					for _, line := range lines {
						result = append(result, line)
					}

					return strings.Join(result, "\n"), nil

				} else {
					tool.mutex.Unlock()
					return "", fmt.Errorf("findings.List: %s", err0.Error())
				}

			} else {
				return "", fmt.Errorf("findings.List: Invalid parameter \"confidence\".")
			}
		} else {
			return "", fmt.Errorf("findings.List: Invalid parameter \"status\".")
		}
	} else {
		return "", fmt.Errorf("findings.List: Invalid parameter \"severity\".")
	}
}
