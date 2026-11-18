package tools

import utils_findings "exocomp/utils/findings"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestFindings_ExportFindingAndEvidence(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)

	_, err0 := tool.BindEvidence(id, findingArguments{
		"kind":    "command_output",
		"content": "HTTP/1.1 200 OK",
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if len(finding.Evidence) != 1 {
		t.Fatalf("Expected one evidence entry, got %d", len(finding.Evidence))
	}

	finding_result, err2 := tool.Export(findingArguments{"type": "finding", "id": id})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err2, finding_result)
	}

	if strings.Contains(finding_result, "./work/findings/") == false {
		t.Errorf("Expected finding export under ./work, got %s", finding_result)
	}

	if strings.Contains(finding_result, ".exocomp") == true {
		t.Errorf("Expected export result not to leak internal state, got %s", finding_result)
	}

	finding_path := filepath.Join(sandbox, "work", "findings", id+".json")

	if _, err3 := os.Stat(finding_path); err3 != nil {
		t.Errorf("Expected exported finding file: %v", err3)
	}

	evidence_result, err4 := tool.Export(findingArguments{
		"type":        "evidence",
		"id":          id,
		"evidence_id": finding.Evidence[0].ID,
	})

	if err4 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err4, evidence_result)
	}

	if strings.Contains(evidence_result, "./work/evidence/") == false {
		t.Errorf("Expected evidence export under ./work, got %s", evidence_result)
	}

	matches, err5 := filepath.Glob(filepath.Join(sandbox, "work", "evidence", id, "*"))

	if err5 != nil {
		t.Fatalf("Expected %v to be nil", err5)
	}

	if len(matches) != 1 {
		t.Fatalf("Expected one exported evidence file, got %d", len(matches))
	}

	content, err6 := os.ReadFile(matches[0])

	if err6 != nil {
		t.Fatalf("Expected %v to be nil", err6)
	}

	if strings.Contains(string(content), "HTTP/1.1 200 OK") == false {
		t.Errorf("Expected exported evidence content, got %q", string(content))
	}

}

func TestFindings_ExportReport(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)

	result, err0 := tool.Export(findingArguments{"type": "report"})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "./work/findings/report.md") == false {
		t.Errorf("Expected report export under ./work, got %s", result)
	}

	if strings.Contains(result, ".exocomp") == true {
		t.Errorf("Expected report export not to leak internal state, got %s", result)
	}

	content, err1 := os.ReadFile(filepath.Join(sandbox, "work", "findings", "report.md"))

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if strings.Contains(string(content), id) == false {
		t.Errorf("Expected report to contain finding ID")
	}

	if strings.Contains(string(content), ".exocomp") == true {
		t.Errorf("Expected report content not to leak internal state")
	}

}

func TestFindings_ExportScan(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	path := writeTestArtifact(t, sandbox, testNucleiResult+"\n")

	result, err0 := tool.Import(findingArguments{"path": path})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	scan_id := ""

	if index := strings.Index(result, "scan "); index >= 0 {
		scan_id = strings.TrimSuffix(result[index+len("scan "):], ").")
	}

	if scan_id == "" {
		t.Fatalf("Expected scan ID in %q", result)
	}

	exported, err1 := tool.Export(findingArguments{"type": "scan", "id": scan_id})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err1, exported)
	}

	if strings.Contains(exported, "./work/scans/") == false {
		t.Errorf("Expected scan export under ./work, got %s", exported)
	}

	content, err2 := os.ReadFile(filepath.Join(sandbox, "work", "scans", scan_id+"-results.jsonl"))

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	if strings.Contains(string(content), "CVE-2024-1234") == false {
		t.Errorf("Expected exported scan content")
	}

}

func TestFindings_GetHidesInternalPaths(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	id := addTestFinding(t, tool)

	proof := filepath.Join(sandbox, "proof_test.go")

	if err := os.WriteFile(proof, []byte("package proof"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	_, err0 := tool.BindEvidence(id, findingArguments{
		"kind":     "unit_test",
		"artifact": "./proof_test.go",
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	result, err1 := tool.Get(id)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if strings.Contains(result, ".exocomp") == true {
		t.Errorf("Expected findings.Get not to leak internal state, got %s", result)
	}

	if strings.Contains(result, "\"artifact\":") == true {
		t.Errorf("Expected findings.Get to omit artifact paths, got %s", result)
	}

}

func TestFindings_ExportRejectsInternalState(t *testing.T) {

	tool, _ := newTestFindings(t)
	_ = addTestFinding(t, tool)

	_, err0 := tool.Export(findingArguments{"type": "report", "path": ".exocomp/report.md"})

	if err0 == nil {
		t.Fatal("Expected export into .exocomp to be rejected")
	}

}
