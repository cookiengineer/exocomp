package findings

import "exocomp/types"
import "testing"

func TestNewFindingID(t *testing.T) {

	id := NewFindingID()

	if types.IsFindingID(id) == false {
		t.Errorf("Expected valid finding ID, got %q", id)
	}

}

func TestNewEvidenceID(t *testing.T) {

	id := NewEvidenceID()

	if len(id) != 11 || id[0:3] != "ev_" {
		t.Errorf("Expected evidence ID, got %q", id)
	}

}
