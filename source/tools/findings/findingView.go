package findings

import "exocomp/schemas"
import "exocomp/types"

type findingView struct {
	ID              string                  `json:"id"`
	Title           string                  `json:"title"`
	Severity        types.FindingSeverity   `json:"severity"`
	CVSS            float64                 `json:"cvss"`
	CWE             string                  `json:"cwe"`
	OWASP           string                  `json:"owasp"`
	AssetIDs        []string                `json:"asset_ids"`
	Endpoint        string                  `json:"endpoint"`
	Parameter       string                  `json:"parameter"`
	Status          types.FindingStatus     `json:"status"`
	Confidence      types.FindingConfidence `json:"confidence"`
	ProofBasis      types.FindingProofBasis `json:"proof_basis"`
	Preconditions   string                  `json:"preconditions"`
	Evidence        []findingEvidenceView   `json:"evidence"`
	Reproduction    string                  `json:"reproduction"`
	Impact          string                  `json:"impact"`
	Remediation     string                  `json:"remediation"`
	Counterevidence string                  `json:"counterevidence"`
	SourceScanID    string                  `json:"source_scan_id"`
	CreatedAt       schemas.Datetime        `json:"created_at"`
	UpdatedAt       schemas.Datetime        `json:"updated_at"`
	VerifiedAt      schemas.Datetime        `json:"verified_at,omitempty"`
	VerifiedBy      string                  `json:"verified_by,omitempty"`
}

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
