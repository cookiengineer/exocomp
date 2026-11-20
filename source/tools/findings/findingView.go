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
