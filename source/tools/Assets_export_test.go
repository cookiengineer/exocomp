package tools

import utils_chat "exocomp/utils/chat"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestAssets_ExportHosts(t *testing.T) {

	tool, sandbox := newTestAssets(t)

	_, err0 := tool.Add(utils_chat.Arguments{
		"type":     "subdomain",
		"value":    "api.example.com",
		"in_scope": true,
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	_, err1 := tool.Add(utils_chat.Arguments{
		"type":  "subdomain",
		"value": "dev.example.com",
	})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	result, err2 := tool.Export(utils_chat.Arguments{
		"format":   "hosts",
		"path":     "./work/targets.txt",
		"in_scope": true,
	})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err2, result)
	}

	if strings.Contains(result, "./work/targets.txt") == false {
		t.Errorf("Expected export path, got %s", result)
	}

	if strings.Contains(result, ".exocomp") == true {
		t.Errorf("Expected export result not to leak internal state, got %s", result)
	}

	content, err3 := os.ReadFile(filepath.Join(sandbox, "work", "targets.txt"))

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

	if string(content) != "api.example.com\n" {
		t.Errorf("Expected only in-scope host, got %q", string(content))
	}

}

func TestAssets_ExportURLsAndJSON(t *testing.T) {

	tool, sandbox := newTestAssets(t)

	_, _ = tool.Add(utils_chat.Arguments{"type": "url", "value": "https://example.com/admin"})
	_, _ = tool.Add(utils_chat.Arguments{"type": "subdomain", "value": "api.example.com"})

	result1, err0 := tool.Export(utils_chat.Arguments{"format": "urls", "path": "./work/urls.txt"})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result1)
	}

	content1, err1 := os.ReadFile(filepath.Join(sandbox, "work", "urls.txt"))

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if string(content1) != "https://example.com/admin\n" {
		t.Errorf("Expected only URL assets, got %q", string(content1))
	}

	result2, err2 := tool.Export(utils_chat.Arguments{"format": "json", "path": "./work/assets.json"})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err2, result2)
	}

	content2, err3 := os.ReadFile(filepath.Join(sandbox, "work", "assets.json"))

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

	if strings.Contains(string(content2), "api.example.com") == false {
		t.Errorf("Expected JSON export to contain assets")
	}

}

func TestAssets_ExportRejectsInternalState(t *testing.T) {

	tool, _ := newTestAssets(t)

	_, err0 := tool.Export(utils_chat.Arguments{"format": "hosts", "path": ".exocomp/targets.txt"})

	if err0 == nil {
		t.Fatal("Expected export into .exocomp to be rejected")
	}

}
