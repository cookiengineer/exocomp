package findings

import "fmt"
import "path/filepath"

func exportReportData(tool *Findings) ([]byte, string, string, error) {

	report, _, err0 := buildFindingsReport(tool, "", false, "", false)

	if err0 == nil {
		return []byte(report), filepath.Join("work", "findings", "report.md"), "text/markdown", nil
	} else {
		return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
	}

}
