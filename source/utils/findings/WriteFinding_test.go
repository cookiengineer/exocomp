package findings

import "exocomp/schemas"
import "exocomp/types"
import "os"
import "path/filepath"
import "testing"

func newFindingsTestPlayground(t *testing.T) string {

	t.Helper()

	playground, err0 := os.MkdirTemp("/tmp", "exocomp-test-utils-findings-*")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	t.Cleanup(func() {
		os.RemoveAll(playground)
	})

	return playground

}

func createTestFinding() types.Finding {

	now := schemas.NewDatetime()

	return types.Finding{
		ID:         "fnd_4a91c2e0",
		Title:      "Authentication bypass",
		Severity:   types.FindingSeverityHigh,
		Status:     types.FindingStatusCandidate,
		Confidence: types.FindingConfidenceUnverified,
		ProofBasis: types.FindingProofBasisStaticUnexecuted,
		CreatedAt:  now,
		UpdatedAt:  now,
		Evidence: []types.FindingEvidence{
			{
				ID:      "ev_4a91c2e0",
				Kind:    types.FindingEvidenceKindNote,
				Title:   "Observation",
				Content: "observed",
			},
		},
	}

}

func TestWriteFinding_RoundTrip(t *testing.T) {

	playground := newFindingsTestPlayground(t)
	finding := createTestFinding()

	err0 := WriteFinding(playground, finding)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	result, err1 := ReadFinding(playground, finding.ID)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if result.Title != finding.Title {
		t.Errorf("Expected title %q, got %q", finding.Title, result.Title)
	}

}

func TestReadFinding_InvalidID(t *testing.T) {

	playground := newFindingsTestPlayground(t)

	_, err0 := ReadFinding(playground, "../index")

	if err0 == nil {
		t.Errorf("Expected invalid finding ID to be refused")
	}

}

func TestWriteIndex_ReadIndex(t *testing.T) {

	playground := newFindingsTestPlayground(t)
	finding := createTestFinding()
	summary := types.FindingSummaryFromFinding(finding)

	err0 := WriteIndex(playground, []types.FindingSummary{summary})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	summaries, err1 := ReadIndex(playground)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if len(summaries) != 1 || summaries[0].ID != finding.ID {
		t.Errorf("Expected one summary for %s, got %v", finding.ID, summaries)
	}

	index_path := filepath.Join(playground, ".exocomp", "findings", "index.json")

	if _, err2 := os.Stat(index_path); err2 != nil {
		t.Errorf("Expected findings index to exist")
	}

}

func TestReadIndex_MissingFile(t *testing.T) {

	playground := newFindingsTestPlayground(t)

	summaries, err0 := ReadIndex(playground)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	} else if len(summaries) != 0 {
		t.Errorf("Expected empty index, got %v", summaries)
	}

}

func TestReadAll(t *testing.T) {

	playground := newFindingsTestPlayground(t)
	finding := createTestFinding()

	err0 := WriteFinding(playground, finding)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	err1 := WriteIndex(playground, []types.FindingSummary{types.FindingSummaryFromFinding(finding)})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	findings, err2 := ReadAll(playground)

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	} else if len(findings) != 1 {
		t.Errorf("Expected one finding, got %d", len(findings))
	}

}
