package types

import "exocomp/schemas"
import "testing"

func TestFinding_Validate(t *testing.T) {

	now := schemas.NewDatetime()

	finding := Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   FindingSeverityHigh,
		Status:     FindingStatusCandidate,
		Confidence: FindingConfidenceUnverified,
		ProofBasis: FindingProofBasisStaticUnexecuted,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if finding.Validate() == true {
		// expected
	} else {
		t.Errorf("Expected finding to be valid")
	}

}

func TestFinding_ValidateInvalidEnums(t *testing.T) {

	now := schemas.NewDatetime()

	finding := Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   FindingSeverity("urgent"),
		Status:     FindingStatusCandidate,
		Confidence: FindingConfidenceUnverified,
		ProofBasis: FindingProofBasisStaticUnexecuted,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if finding.Validate() == false {
		// expected
	} else {
		t.Errorf("Expected invalid severity to be rejected")
	}

	finding.Severity = FindingSeverityHigh
	finding.Status = FindingStatus("unknown")

	if finding.Validate() == false {
		// expected
	} else {
		t.Errorf("Expected invalid status to be rejected")
	}

}

func TestFinding_ValidateCVSSBounds(t *testing.T) {

	now := schemas.NewDatetime()

	finding := Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   FindingSeverityHigh,
		CVSS:       11,
		Status:     FindingStatusCandidate,
		Confidence: FindingConfidenceUnverified,
		ProofBasis: FindingProofBasisStaticUnexecuted,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if finding.Validate() == false {
		// expected
	} else {
		t.Errorf("Expected CVSS above 10 to be rejected")
	}

}

func TestFinding_ValidateEvidence(t *testing.T) {

	now := schemas.NewDatetime()

	finding := Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   FindingSeverityHigh,
		Status:     FindingStatusCandidate,
		Confidence: FindingConfidenceUnverified,
		ProofBasis: FindingProofBasisStaticUnexecuted,
		CreatedAt:  now,
		UpdatedAt:  now,
		Evidence: []FindingEvidence{
			{
				ID:   "ev_12345678",
				Kind: FindingEvidenceKindNote,
			},
		},
	}

	if finding.Validate() == false {
		// expected
	} else {
		t.Errorf("Expected evidence without content or artifact to be rejected")
	}

	finding.Evidence[0].Content = "observed"

	if finding.Validate() == true {
		// expected
	} else {
		t.Errorf("Expected evidence with content to be valid")
	}

}

func TestFindingSummaryFromFinding(t *testing.T) {

	now := schemas.NewDatetime()

	finding := Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   FindingSeverityCritical,
		Status:     FindingStatusVerified,
		Confidence: FindingConfidenceObserved,
		ProofBasis: FindingProofBasisExecutedLab,
		AssetIDs:   []string{"asset-1"},
		CreatedAt:  now,
		UpdatedAt:  now,
		Evidence: []FindingEvidence{
			{ID: "ev_12345678", Kind: FindingEvidenceKindNote, Content: "observed"},
		},
	}

	summary := FindingSummaryFromFinding(finding)

	if summary.ID == finding.ID && summary.EvidenceCount == 1 && len(summary.AssetIDs) == 1 {
		// expected
	} else {
		t.Errorf("Expected summary to reflect finding")
	}

}

func TestIsFindingID(t *testing.T) {

	if IsFindingID("fnd_4a91c2e0") == true {
		// expected
	} else {
		t.Errorf("Expected valid finding ID to be accepted")
	}

	if IsFindingID("fnd_4a91c2ez") == false {
		// expected
	} else {
		t.Errorf("Expected invalid finding ID to be rejected")
	}

	if IsFindingID("../fnd_4a91c2e0") == false {
		// expected
	} else {
		t.Errorf("Expected path traversal finding ID to be rejected")
	}

}
