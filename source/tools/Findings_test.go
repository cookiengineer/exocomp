package tools

import utils_findings "exocomp/utils/findings"
import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "os"
import "path/filepath"
import "strings"
import "testing"

func newTestFindings(t *testing.T) (*Findings, string) {

	t.Helper()

	playground, err0 := os.MkdirTemp("/tmp", "exocomp-test-findings-*")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	sandbox := filepath.Join(playground, "sandbox")

	err1 := os.MkdirAll(sandbox, 0755)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

	tool := NewFindings([]string{
		"Add", "BindEvidence", "Export", "Get", "Import", "List", "Report", "Retest", "Search", "Update", "Verify",
	}, playground, sandbox, "Peanut Reproducer")

	return tool, sandbox

}

func addTestFinding(t *testing.T, tool *Findings) string {

	t.Helper()

	result, err := tool.Add(findingArguments{
		"title":    "Authentication bypass",
		"severity": "high",
		"endpoint": "https://example.com/admin",
		"cvss":     float64(8.1),
	})

	if err != nil {
		t.Fatalf("Expected %v to be nil, got %s", err, result)
	}

	ids := strings.Split(result, "\"")

	if len(ids) >= 2 {
		return ids[1]
	}

	t.Fatalf("Expected finding ID in %q", result)

	return ""

}

func TestFindings_BindEvidenceFromExocomp(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)

	path := filepath.Join(sandbox, ".exocomp", "scans", "evidence.txt")

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(path, []byte("scanner evidence"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	result, err0 := tool.BindEvidence(id, findingArguments{
		"kind":     "artifact",
		"artifact": ".exocomp/scans/evidence.txt",
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if len(finding.Evidence) != 1 {
		t.Fatalf("Expected one evidence entry, got %d", len(finding.Evidence))
	}

	if strings.HasPrefix(finding.Evidence[0].Artifact, "./.exocomp/evidence/") == false {
		t.Errorf("Expected evidence artifact in .exocomp/evidence, got %q", finding.Evidence[0].Artifact)
	}

}

func TestFindings_AddAndGet(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	} else if finding.Title != "Authentication bypass" {
		t.Errorf("Expected finding title, got %q", finding.Title)
	}

	result, err1 := tool.Get(id)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if strings.Contains(result, "\"status\": \"candidate\"") == false {
		t.Errorf("Expected candidate status in %s", result)
	}

}

func TestFindings_ListSearchUpdate(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	_, err0 := tool.Add(findingArguments{
		"title":    "Information disclosure",
		"severity": "low",
		"endpoint": "https://example.com/debug",
	})

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	result1, err1 := tool.List(findingArguments{"severity": "high"})

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if strings.Contains(result1, "1 findings") == false {
		t.Errorf("Expected one high finding, got %s", result1)
	}

	result2, err2 := tool.Search(findingArguments{"query": "debug"})

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	} else if strings.Contains(result2, "Information disclosure") == false {
		t.Errorf("Expected search to find disclosure, got %s", result2)
	}

	result3, err3 := tool.Update(id, findingArguments{"title": "Authentication bypass confirmed"})

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	} else if strings.Contains(result3, "updated") == false {
		t.Errorf("Expected update confirmation, got %s", result3)
	}

	_, err4 := tool.Update(id, findingArguments{"status": "verified"})

	if err4 == nil {
		t.Errorf("Expected status update through findings.Update to fail")
	}

}

func TestFindings_BindEvidenceContent(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	result, err0 := tool.BindEvidence(id, findingArguments{
		"kind":    "http_response",
		"content": "Authorization: Bearer supersecret\nHTTP/1.1 200 OK",
	})

	if err0 != nil {
		t.Errorf("Expected %v to be nil, got %s", err0, result)
	}

	finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	} else if len(finding.Evidence) != 1 {
		t.Fatalf("Expected one evidence entry")
	}

	evidence := finding.Evidence[0]

	if strings.Contains(evidence.Content, "supersecret") == true {
		t.Errorf("Expected bearer token to be redacted")
	}

	sum := sha256.Sum256([]byte(evidence.Content))
	expected_hash := hex.EncodeToString(sum[:])

	if evidence.SHA256 != expected_hash {
		t.Errorf("Expected SHA-256 %s, got %s", expected_hash, evidence.SHA256)
	}

	if strings.HasPrefix(evidence.Artifact, "./.exocomp/evidence/") == false {
		t.Errorf("Expected relative evidence artifact path, got %s", evidence.Artifact)
	}

}

func TestFindings_BindEvidenceArtifact(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)
	proof_path := filepath.Join(sandbox, "proof_test.go")
	proof := "package findings\n\nfunc TestProof(t *testing.T) {}\n"

	err0 := os.WriteFile(proof_path, []byte(proof), 0644)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	_, err1 := tool.BindEvidence(id, findingArguments{
		"kind":     "unit_test",
		"artifact": "./proof_test.go",
	})

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	finding, err2 := utils_findings.ReadFinding(tool.Playground, id)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	sum := sha256.Sum256([]byte(proof))
	expected_hash := hex.EncodeToString(sum[:])

	if finding.Evidence[0].SHA256 != expected_hash {
		t.Errorf("Expected SHA-256 %s, got %s", expected_hash, finding.Evidence[0].SHA256)
	}

}

func TestFindings_VerifyRequiresEvidence(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	_, err0 := tool.Verify(id, findingArguments{"proof_basis": "executed_lab"})

	if err0 == nil {
		t.Errorf("Expected verification without evidence to fail")
	}

}

func TestFindings_VerifyRequiresReproduction(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)
	proof_path := filepath.Join(sandbox, "proof_test.go")

	err0 := os.WriteFile(proof_path, []byte("package findings\n"), 0644)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	_, err1 := tool.BindEvidence(id, findingArguments{
		"kind":     "unit_test",
		"artifact": "./proof_test.go",
	})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	_, err2 := tool.BindEvidence(id, findingArguments{
		"kind":        "test_run",
		"command":     "go test ./findings/fnd_test",
		"exit_code":   0,
		"content":     "PASS",
		"environment": map[string]string{"EXOCOMP_TARGET": "https://example.com"},
	})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	_, err3 := tool.Verify(id, findingArguments{
		"proof_basis":          "executed_lab",
		"require_reproduction": true,
	})

	if err3 != nil {
		t.Errorf("Expected %v to be nil, got %s", err3, err3.Error())
	}

	finding, err4 := utils_findings.ReadFinding(tool.Playground, id)

	if err4 != nil {
		t.Fatalf("Expected %v to be nil", err4)
	} else if finding.Status != "verified" {
		t.Errorf("Expected verified status, got %s", finding.Status)
	} else if finding.VerifiedBy != "Peanut Reproducer" {
		t.Errorf("Expected verifier name, got %s", finding.VerifiedBy)
	}

}

func TestFindings_Retest(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	_, err0 := tool.BindEvidence(id, findingArguments{
		"kind":    "note",
		"content": "Fixed in version 2.0.",
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	_, err1 := tool.Retest(id, findingArguments{
		"status":          "fixed",
		"counterevidence": "Version 2.0 returns 403.",
	})

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	finding, err2 := utils_findings.ReadFinding(tool.Playground, id)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	} else if finding.Status != "fixed" {
		t.Errorf("Expected fixed status, got %s", finding.Status)
	}

}

func TestFindings_Report(t *testing.T) {

	tool, _ := newTestFindings(t)
	id := addTestFinding(t, tool)

	_, err0 := tool.BindEvidence(id, findingArguments{
		"kind":    "command_output",
		"content": "HTTP/1.1 200 OK",
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	result, err1 := tool.Report(findingArguments{})

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if strings.Contains(result, "Authentication bypass") == false {
		t.Errorf("Expected report content, got %s", result)
	} else if strings.Contains(result, ".exocomp") == true {
		t.Errorf("Expected report content not to leak internal state, got %s", result)
	}

	report_path := filepath.Join(tool.Playground, ".exocomp", "findings", "report.md")
	report, err2 := os.ReadFile(report_path)

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	} else if strings.Contains(string(report), "Authentication bypass") == false {
		t.Errorf("Expected report to contain finding title")
	} else if strings.Contains(string(report), "SHA-256:") == false {
		t.Errorf("Expected report to contain evidence hash")
	}

}

func TestFindings_InvalidID(t *testing.T) {

	tool, _ := newTestFindings(t)

	_, err0 := tool.Get("../index")

	if err0 == nil {
		t.Errorf("Expected invalid ID to be refused")
	}

	_, err1 := tool.BindEvidence("../index", findingArguments{
		"kind":    "note",
		"content": "test",
	})

	if err1 == nil {
		t.Errorf("Expected invalid ID to be refused")
	}

}

func TestFindings_SchemaMatchesMethods(t *testing.T) {

	tool, _ := newTestFindings(t)
	schemas := tool.Schemas()

	if len(schemas) != len(tool.Methods) {
		t.Errorf("Expected %d schemas, got %d", len(tool.Methods), len(schemas))
	}

}

func TestFindings_MarshalJSON(t *testing.T) {

	tool, _ := newTestFindings(t)
	addTestFinding(t, tool)

	buffer, err0 := json.Marshal(tool)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	} else if strings.Contains(string(buffer), "Authentication bypass") == false {
		t.Errorf("Expected index JSON to contain finding")
	}

}
