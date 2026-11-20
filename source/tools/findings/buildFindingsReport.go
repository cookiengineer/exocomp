package findings

import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import "strings"
import "time"

func buildFindingsReport(tool *Findings, severity_filter string, has_severity bool, status_filter string, has_status bool) (string, int, error) {

	tool.mutex.Lock()
	findings, err0 := utils_findings.ReadAll(tool.Playground)
	tool.mutex.Unlock()

	if err0 == nil {

		filtered := make([]types.Finding, 0)

		for _, finding := range findings {

			if has_severity == true && string(finding.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
				continue
			} else if has_status == true && string(finding.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
				continue
			}

			filtered = append(filtered, finding)

		}

		report := utils_findings.BuildReport(filtered, time.Now().Format("2006-01-02 15:04:05"))

		if err1 := writeFindingsReport(tool, report); err1 == nil {
			return report, len(filtered), nil
		} else {
			return "", 0, err1
		}

	} else {
		return "", 0, err0
	}

}
