package tools

import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"
import "path/filepath"
import "strings"
import "time"

func buildFindingsReport(tool *Findings, severity_filter string, has_severity bool, status_filter string, has_status bool) (string, int, error) {

	tool.mutex.Lock()
	findings, err0 := utils_findings.ReadAll(tool.Playground)
	tool.mutex.Unlock()

	if err0 != nil {
		return "", 0, err0
	}

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
	resolved, err1 := utils_sandbox.ResolvePath(tool.Playground, filepath.Join(".exocomp", "findings", "report.md"))

	if err1 != nil {
		return "", 0, fmt.Errorf("Cannot resolve findings report.")
	}

	if err2 := os.WriteFile(resolved, []byte(report), 0666); err2 != nil {
		return "", 0, fmt.Errorf("Cannot write findings report.")
	}

	return report, len(filtered), nil

}
