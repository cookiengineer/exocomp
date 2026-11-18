package findings

import "exocomp/types"
import "testing"

func TestSeverityRank(t *testing.T) {

	if SeverityRank(types.FindingSeverityCritical) != 0 {
		t.Errorf("Expected critical rank 0")
	}

	if SeverityRank(types.FindingSeverityInfo) != 4 {
		t.Errorf("Expected info rank 4")
	}

	if SeverityRank(types.FindingSeverity("unknown")) != 5 {
		t.Errorf("Expected unknown rank 5")
	}

}
