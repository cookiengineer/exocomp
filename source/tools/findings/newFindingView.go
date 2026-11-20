package findings

import "exocomp/types"

func newFindingView(finding types.Finding) findingView {

	evidence := make([]findingEvidenceView, 0, len(finding.Evidence))

	for _, entry := range finding.Evidence {
		evidence = append(evidence, findingEvidenceView{
			ID:          entry.ID,
			Kind:        entry.Kind,
			Title:       entry.Title,
			Source:      entry.Source,
			Command:     entry.Command,
			ExitCode:    entry.ExitCode,
			Environment: entry.Environment,
			Content:     entry.Content,
			SHA256:      entry.SHA256,
			Size:        entry.Size,
			ContentType: entry.ContentType,
			CreatedAt:   entry.CreatedAt,
		})
	}

	return findingView{
		ID:              finding.ID,
		Title:           finding.Title,
		Severity:        finding.Severity,
		CVSS:            finding.CVSS,
		CWE:             finding.CWE,
		OWASP:           finding.OWASP,
		AssetIDs:        finding.AssetIDs,
		Endpoint:        finding.Endpoint,
		Parameter:       finding.Parameter,
		Status:          finding.Status,
		Confidence:      finding.Confidence,
		ProofBasis:      finding.ProofBasis,
		Preconditions:   finding.Preconditions,
		Evidence:        evidence,
		Reproduction:    finding.Reproduction,
		Impact:          finding.Impact,
		Remediation:     finding.Remediation,
		Counterevidence: finding.Counterevidence,
		SourceScanID:    finding.SourceScanID,
		CreatedAt:       finding.CreatedAt,
		UpdatedAt:       finding.UpdatedAt,
		VerifiedAt:      finding.VerifiedAt,
		VerifiedBy:      finding.VerifiedBy,
	}

}
