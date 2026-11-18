package findings

import "exocomp/types"
import "testing"

func TestMatchesQuery(t *testing.T) {

	finding := types.Finding{
		ID:       "fnd_4a91c2e0",
		Title:    "Authentication bypass",
		Endpoint: "https://example.com/admin",
		Evidence: []types.FindingEvidence{
			{
				ID:      "ev_4a91c2e0",
				Kind:    types.FindingEvidenceKindNote,
				Title:   "Response",
				Content: "observed bypass",
			},
		},
	}

	if MatchesQuery(finding, "authentication") == false {
		t.Errorf("Expected title match")
	}

	if MatchesQuery(finding, "example.com") == false {
		t.Errorf("Expected endpoint match")
	}

	if MatchesQuery(finding, "bypass") == false {
		t.Errorf("Expected evidence match")
	}

	if MatchesQuery(finding, "missing") == true {
		t.Errorf("Expected missing query not to match")
	}

}
