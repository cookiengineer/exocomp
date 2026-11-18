package findings

import "exocomp/types"
import "fmt"
import "sort"
import "strings"

func BuildReport(findings []types.Finding, generated string) string {

	sort.Slice(findings, func(a int, b int) bool {

		rank1 := SeverityRank(findings[a].Severity)
		rank2 := SeverityRank(findings[b].Severity)

		if rank1 != rank2 {
			return rank1 < rank2
		} else if findings[a].UpdatedAt.Equal(findings[b].UpdatedAt.Time) == false {
			return findings[a].UpdatedAt.After(findings[b].UpdatedAt.Time)
		} else {
			return findings[a].ID < findings[b].ID
		}

	})

	severity_counts := make(map[types.FindingSeverity]int, 0)
	status_counts := make(map[types.FindingStatus]int, 0)

	for _, finding := range findings {
		severity_counts[finding.Severity] = severity_counts[finding.Severity] + 1
		status_counts[finding.Status] = status_counts[finding.Status] + 1
	}

	builder := strings.Builder{}
	builder.WriteString("# Findings Report\n\n")
	builder.WriteString(fmt.Sprintf("Generated: %s\n\n", generated))
	builder.WriteString(fmt.Sprintf("Total findings: %d\n\n", len(findings)))
	builder.WriteString("## Summary\n\n")

	for _, severity := range []types.FindingSeverity{
		types.FindingSeverityCritical,
		types.FindingSeverityHigh,
		types.FindingSeverityMedium,
		types.FindingSeverityLow,
		types.FindingSeverityInfo,
	} {
		builder.WriteString(fmt.Sprintf("- %s: %d\n", severity, severity_counts[severity]))
	}

	builder.WriteString("\n")

	for _, status := range []types.FindingStatus{
		types.FindingStatusCandidate,
		types.FindingStatusVerified,
		types.FindingStatusFixed,
		types.FindingStatusFalsePositive,
		types.FindingStatusAcceptedRisk,
	} {
		builder.WriteString(fmt.Sprintf("- %s: %d\n", status, status_counts[status]))
	}

	builder.WriteString("\n## Findings\n\n")

	for _, finding := range findings {

		builder.WriteString(fmt.Sprintf("### %s - %s\n\n", finding.ID, finding.Title))
		builder.WriteString(fmt.Sprintf("- Severity: %s\n", finding.Severity))
		builder.WriteString(fmt.Sprintf("- CVSS: %.1f\n", finding.CVSS))
		builder.WriteString(fmt.Sprintf("- Status: %s\n", finding.Status))
		builder.WriteString(fmt.Sprintf("- Confidence: %s\n", finding.Confidence))
		builder.WriteString(fmt.Sprintf("- Proof basis: %s\n", finding.ProofBasis))

		if finding.CWE != "" {
			builder.WriteString(fmt.Sprintf("- CWE: %s\n", finding.CWE))
		}

		if finding.OWASP != "" {
			builder.WriteString(fmt.Sprintf("- OWASP: %s\n", finding.OWASP))
		}

		if finding.Endpoint != "" {
			builder.WriteString(fmt.Sprintf("- Endpoint: %s\n", finding.Endpoint))
		}

		if finding.Parameter != "" {
			builder.WriteString(fmt.Sprintf("- Parameter: %s\n", finding.Parameter))
		}

		if len(finding.AssetIDs) > 0 {
			builder.WriteString(fmt.Sprintf("- Assets: %s\n", strings.Join(finding.AssetIDs, ", ")))
		}

		if finding.VerifiedAt.IsZero() == false {
			builder.WriteString(fmt.Sprintf("- Verified at: %s\n", finding.VerifiedAt.Format("2006-01-02 15:04:05")))
		}

		if finding.VerifiedBy != "" {
			builder.WriteString(fmt.Sprintf("- Verified by: %s\n", finding.VerifiedBy))
		}

		if finding.Preconditions != "" {
			builder.WriteString(fmt.Sprintf("\n#### Preconditions\n\n%s\n", finding.Preconditions))
		}

		if finding.Reproduction != "" {
			builder.WriteString(fmt.Sprintf("\n#### Reproduction\n\n%s\n", finding.Reproduction))
		}

		if finding.Impact != "" {
			builder.WriteString(fmt.Sprintf("\n#### Impact\n\n%s\n", finding.Impact))
		}

		if finding.Remediation != "" {
			builder.WriteString(fmt.Sprintf("\n#### Remediation\n\n%s\n", finding.Remediation))
		}

		if finding.Counterevidence != "" {
			builder.WriteString(fmt.Sprintf("\n#### Counterevidence\n\n%s\n", finding.Counterevidence))
		}

		if len(finding.Evidence) > 0 {

			builder.WriteString("\n#### Evidence\n\n")

			for _, evidence := range finding.Evidence {

				builder.WriteString(fmt.Sprintf("- %s: %s\n", evidence.ID, evidence.Title))

				if evidence.SHA256 != "" {
					builder.WriteString(fmt.Sprintf("  - SHA-256: %s\n", evidence.SHA256))
				}

				if evidence.Size > 0 {
					builder.WriteString(fmt.Sprintf("  - Size: %d B\n", evidence.Size))
				}

				if evidence.Command != "" {
					builder.WriteString(fmt.Sprintf("  - Command: %s\n", evidence.Command))
				}

				if evidence.ExitCode != nil {
					builder.WriteString(fmt.Sprintf("  - Exit code: %d\n", *evidence.ExitCode))
				}

				if len(evidence.Environment) > 0 {

					for key, value := range evidence.Environment {
						builder.WriteString(fmt.Sprintf("  - Environment: %s=%s\n", key, value))
					}

				}

			}

		}

		builder.WriteString("\n")

	}

	return builder.String()

}
