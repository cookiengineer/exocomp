package findings

import "exocomp/types"
import "strings"
import "testing"

func TestBuildReport(t *testing.T) {

	finding := createTestFinding()
	finding.Evidence[0].SHA256 = "0123456789abcdef"
	finding.Evidence[0].Artifact = "./.exocomp/evidence/fnd_4a91c2e0/ev_4a91c2e0-evidence.txt"

	report := BuildReport([]types.Finding{finding}, "2026-11-18 14:00:00")

	if strings.Contains(report, "Authentication bypass") == false {
		t.Errorf("Expected report to contain finding title")
	}

	if strings.Contains(report, "SHA-256: 0123456789abcdef") == false {
		t.Errorf("Expected report to contain evidence hash")
	}

	if strings.Contains(report, "Total findings: 1") == false {
		t.Errorf("Expected report to contain finding count")
	}

}
