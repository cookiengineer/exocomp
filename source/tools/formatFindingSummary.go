package tools

import "exocomp/types"
import "fmt"

func formatFindingSummary(summary types.FindingSummary) string {

	return fmt.Sprintf(
		"- ID: %s, Severity: %s, CVSS: %.1f, Status: %s, Confidence: %s, Endpoint: \"%s\", Title: \"%s\"",
		summary.ID,
		summary.Severity,
		summary.CVSS,
		summary.Status,
		summary.Confidence,
		summary.Endpoint,
		summary.Title,
	)

}
