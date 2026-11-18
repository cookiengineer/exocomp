package types

func FindingSummaryFromFinding(finding Finding) FindingSummary {

	return FindingSummary{
		ID:            finding.ID,
		Title:         finding.Title,
		Severity:      finding.Severity,
		CVSS:          finding.CVSS,
		Status:        finding.Status,
		Confidence:    finding.Confidence,
		ProofBasis:    finding.ProofBasis,
		AssetIDs:      append([]string{}, finding.AssetIDs...),
		Endpoint:      finding.Endpoint,
		Parameter:     finding.Parameter,
		EvidenceCount: len(finding.Evidence),
		CreatedAt:     finding.CreatedAt,
		UpdatedAt:     finding.UpdatedAt,
	}

}
