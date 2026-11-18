package types

import "exocomp/schemas"

type Finding struct {
	ID              string            `json:"id"`
	Title           string            `json:"title"`
	Severity        FindingSeverity   `json:"severity"`
	CVSS            float64           `json:"cvss"`
	CWE             string            `json:"cwe"`
	OWASP           string            `json:"owasp"`
	AssetIDs        []string          `json:"asset_ids"`
	Endpoint        string            `json:"endpoint"`
	Parameter       string            `json:"parameter"`
	Status          FindingStatus     `json:"status"`
	Confidence      FindingConfidence `json:"confidence"`
	ProofBasis      FindingProofBasis `json:"proof_basis"`
	Preconditions   string            `json:"preconditions"`
	Evidence        []FindingEvidence `json:"evidence"`
	Reproduction    string            `json:"reproduction"`
	Impact          string            `json:"impact"`
	Remediation     string            `json:"remediation"`
	Counterevidence string            `json:"counterevidence"`
	SourceScanID    string            `json:"source_scan_id"`
	CreatedAt       schemas.Datetime  `json:"created_at"`
	UpdatedAt       schemas.Datetime  `json:"updated_at"`
	VerifiedAt      schemas.Datetime  `json:"verified_at,omitempty"`
	VerifiedBy      string            `json:"verified_by,omitempty"`
}

func (finding *Finding) Validate() bool {

	if IsFindingID(finding.ID) == false {
		return false
	} else if finding.Title == "" {
		return false
	} else if finding.Severity.IsValid() == false {
		return false
	} else if finding.Status.IsValid() == false {
		return false
	} else if finding.Confidence.IsValid() == false {
		return false
	} else if finding.ProofBasis.IsValid() == false {
		return false
	} else if finding.CVSS < 0 || finding.CVSS > 10 {
		return false
	}

	for _, evidence := range finding.Evidence {

		if evidence.Kind.IsValid() == false {
			return false
		} else if evidence.Content == "" && evidence.Artifact == "" {
			return false
		}

	}

	return true

}
