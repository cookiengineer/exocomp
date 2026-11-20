package types

import "exocomp/schemas"

type FindingSummary struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Severity      FindingSeverity   `json:"severity"`
	CVSS          float64           `json:"cvss"`
	Status        FindingStatus     `json:"status"`
	Confidence    FindingConfidence `json:"confidence"`
	ProofBasis    FindingProofBasis `json:"proof_basis"`
	AssetIDs      []string          `json:"asset_ids"`
	Endpoint      string            `json:"endpoint"`
	Parameter     string            `json:"parameter"`
	EvidenceCount int               `json:"evidence_count"`
	CreatedAt     schemas.Datetime  `json:"created_at"`
	UpdatedAt     schemas.Datetime  `json:"updated_at"`
}

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
