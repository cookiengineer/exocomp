package tools

import utils_findings "exocomp/utils/findings"
import "os"
import "path/filepath"
import "strings"
import "testing"

const testNucleiResult = `{"template-id":"CVE-2024-1234","info":{"name":"Example RCE","severity":"high","description":"Remote code execution","classification":{"cwe-id":["CWE-94"],"cvss-score":9.8}},"type":"http","host":"https://example.com","matched-at":"https://example.com/vuln","request":"GET /vuln HTTP/1.1\r\nHost: example.com\r\n\r\n","response":"HTTP/1.1 200 OK\r\n\r\nok"}`

func writeTestArtifact(t *testing.T, sandbox string, content string) string {

	t.Helper()

	path := filepath.Join(sandbox, "work", "scans", "nuclei", "results.jsonl")

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	return "./work/scans/nuclei/results.jsonl"

}

func TestFindings_ImportNuclei(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	path := writeTestArtifact(t, sandbox, testNucleiResult+"\n")

	result, err0 := tool.Import(findingArguments{"path": path})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "Imported 1 candidates") == false {
		t.Fatalf("Expected one imported candidate, got %q", result)
	}

	summaries, err1 := utils_findings.ReadIndex(tool.Playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	} else if len(summaries) != 1 {
		t.Fatalf("Expected one finding, got %d", len(summaries))
	}

	finding, err2 := utils_findings.ReadFinding(tool.Playground, summaries[0].ID)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	if finding.Title != "Example RCE" {
		t.Errorf("Expected title %q, got %q", "Example RCE", finding.Title)
	}

	if finding.Severity != "high" {
		t.Errorf("Expected high severity, got %q", finding.Severity)
	}

	if finding.CWE != "CWE-94" {
		t.Errorf("Expected CWE-94, got %q", finding.CWE)
	}

	if finding.CVSS != 9.8 {
		t.Errorf("Expected CVSS 9.8, got %v", finding.CVSS)
	}

	if strings.HasPrefix(finding.SourceScanID, "scn_") == false {
		t.Errorf("Expected scan ID, got %q", finding.SourceScanID)
	}

	if len(finding.Evidence) != 2 {
		t.Fatalf("Expected two evidence entries, got %d", len(finding.Evidence))
	}

	if finding.Evidence[0].Kind != "http_request" {
		t.Errorf("Expected http_request evidence, got %q", finding.Evidence[0].Kind)
	}

	if finding.Evidence[1].Kind != "http_response" {
		t.Errorf("Expected http_response evidence, got %q", finding.Evidence[1].Kind)
	}

	second, err3 := tool.Import(findingArguments{"path": path})

	if err3 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err3, second)
	}

	if strings.Contains(second, "skipped 1 duplicates") == false {
		t.Errorf("Expected duplicate to be skipped, got %q", second)
	}

}

func TestFindings_ImportNucleiInvalidLines(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	path := writeTestArtifact(t, sandbox, "not json\n"+testNucleiResult+"\n")

	result, err0 := tool.Import(findingArguments{"path": path})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "Imported 1 candidates") == false {
		t.Errorf("Expected one imported candidate, got %q", result)
	}

	if strings.Contains(result, "1 invalid results") == false {
		t.Errorf("Expected one invalid result, got %q", result)
	}

}

func TestFindings_ImportNucleiEmpty(t *testing.T) {

	tool, sandbox := newTestFindings(t)
	path := writeTestArtifact(t, sandbox, "")

	result, err0 := tool.Import(findingArguments{"path": path})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "Imported 0 candidates") == false {
		t.Errorf("Expected zero imported candidates, got %q", result)
	}

}

func TestFindings_ImportRejectsMissingArtifact(t *testing.T) {

	tool, _ := newTestFindings(t)

	_, err0 := tool.Import(findingArguments{"path": "./work/scans/nuclei/missing.jsonl"})

	if err0 == nil {
		t.Fatal("Expected an error for a missing artifact")
	}

	if strings.Contains(err0.Error(), tool.Sandbox) == true {
		t.Errorf("Expected sandbox path to be hidden, got %q", err0.Error())
	}

}
