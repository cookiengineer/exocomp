package tools

import utils_findings "exocomp/utils/findings"
import "os"
import "path/filepath"
import "strings"
import "testing"

func newTestAssets(t *testing.T) (*Assets, string) {

	t.Helper()

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "sandbox")

	if err := os.MkdirAll(sandbox, 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	tool := NewAssets([]string{
		"Add", "Export", "Import", "Get", "List", "Search", "Update", "Neighbors", "Scope", "Coverage",
	}, playground, sandbox)

	return tool, sandbox

}

func writeAssetFixture(t *testing.T, sandbox string, name string, content string) string {

	t.Helper()

	relative := filepath.Join("work", "scans", name, "results.jsonl")
	path := filepath.Join(sandbox, relative)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	return relative

}

func TestAssets_AddDeduplicates(t *testing.T) {

	tool, _ := newTestAssets(t)

	result1, err0 := tool.Add(assetArguments{
		"type":     "subdomain",
		"value":    "API.Example.com",
		"source":   "manual",
		"in_scope": true,
		"tags":     []interface{}{"api"},
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result1)
	}

	result2, err1 := tool.Add(assetArguments{
		"type":  "subdomain",
		"value": "api.example.com",
	})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err1, result2)
	}

	ids := tool.GetContentIdentifiers()

	if len(ids) != 1 {
		t.Fatalf("Expected one deduplicated asset, got %d", len(ids))
	}

	entries, err2 := utils_findings.ReadAssetsIndex(tool.Playground)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	if entries[0].Value != "api.example.com" {
		t.Errorf("Expected normalized value, got %q", entries[0].Value)
	}

	if entries[0].InScope != true {
		t.Errorf("Expected asset to stay in scope")
	}

	scope, err3 := tool.Scope(assetArguments{})

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

	if strings.Contains(scope, "1 assets, 1 in scope") == false {
		t.Errorf("Expected scope summary, got %q", scope)
	}

	coverage, err4 := tool.Coverage()

	if err4 != nil {
		t.Fatalf("Expected %v to be nil", err4)
	}

	if strings.Contains(coverage, "subdomain: 1 assets") == false {
		t.Errorf("Expected subdomain coverage, got %q", coverage)
	}

}

func TestAssets_ImportHttpx(t *testing.T) {

	tool, sandbox := newTestAssets(t)

	path := writeAssetFixture(t, sandbox, "httpx", strings.Join([]string{
		`{"url":"https://api.example.com/v1","host":"api.example.com","status_code":200,"title":"API","tech":["nginx"],"webserver":"nginx","ip":"1.2.3.4"}`,
		`{"url":"https://api.example.com/v1","host":"api.example.com","status_code":200,"title":"API","tech":["nginx"],"webserver":"nginx","ip":"1.2.3.4"}`,
	}, "\n"))

	result, err0 := tool.Import(assetArguments{"path": path, "format": "httpx", "in_scope": true})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "Created 1 assets, updated 1") == false {
		t.Errorf("Expected one created and one updated asset, got %q", result)
	}

	entries, err1 := utils_findings.ReadAssetsIndex(tool.Playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected one asset, got %d", len(entries))
	}

	if entries[0].Type != "url" {
		t.Errorf("Expected URL asset, got %q", entries[0].Type)
	}

	if entries[0].Props["status_code"] != "200" {
		t.Errorf("Expected status code property, got %q", entries[0].Props["status_code"])
	}

	if len(entries[0].Technology) != 1 || entries[0].Technology[0] != "nginx" {
		t.Errorf("Expected nginx technology, got %v", entries[0].Technology)
	}

}

func TestAssets_ImportSubfinderAndScope(t *testing.T) {

	tool, sandbox := newTestAssets(t)
	path := writeAssetFixture(t, sandbox, "subfinder", strings.Join([]string{
		`{"host":"api.example.com","source":"crtsh"}`,
		`{"host":"example.com","source":"crtsh"}`,
	}, "\n"))

	result, err0 := tool.Import(assetArguments{"path": path, "format": "subfinder", "target": "example.com"})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	entries, err1 := utils_findings.ReadAssetsIndex(tool.Playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected two assets, got %d", len(entries))
	}

	types := map[string]int{}

	for _, entry := range entries {
		types[string(entry.Type)]++
	}

	if types["subdomain"] != 1 || types["domain"] != 1 {
		t.Errorf("Expected one subdomain and one domain, got %v", types)
	}

	scoped, err2 := tool.Scope(assetArguments{"type": "subdomain"})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	if strings.Contains(scoped, "1 assets marked in_scope=true") == false {
		t.Errorf("Expected subdomain to be marked in scope, got %q", scoped)
	}

}

func TestAssets_ImportNaabuKatanaGobuster(t *testing.T) {

	tool, sandbox := newTestAssets(t)

	naabu_path := writeAssetFixture(t, sandbox, "naabu", `{"host":"api.example.com","ip":"1.2.3.4","port":443,"protocol":"tcp","tls":true}`)
	katana_path := writeAssetFixture(t, sandbox, "katana", `{"request":{"method":"GET","endpoint":"https://example.com/admin"},"response":{"status_code":200}}`)
	gobuster_path := writeAssetFixture(t, sandbox, "gobuster", "/backup (Status: 301) [Size: 0]\nhttps://example.com/login (Status: 200) [Size: 123]")

	if _, err := tool.Import(assetArguments{"path": naabu_path, "format": "naabu"}); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if _, err := tool.Import(assetArguments{"path": katana_path, "format": "katana"}); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	result, err := tool.Import(assetArguments{"path": gobuster_path, "format": "gobuster", "target": "https://example.com"})

	if err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if strings.Contains(result, "Created 2 assets") == false {
		t.Errorf("Expected two gobuster endpoints, got %q", result)
	}

	entries, err1 := utils_findings.ReadAssetsIndex(tool.Playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	types := map[string]int{}

	for _, entry := range entries {
		types[string(entry.Type)]++
	}

	if types["port"] != 1 {
		t.Errorf("Expected one port asset, got %v", types)
	}

	if types["endpoint"] != 3 {
		t.Errorf("Expected three endpoint assets, got %v", types)
	}

}

func TestAssets_ImportDnsxAndTlsx(t *testing.T) {

	tool, sandbox := newTestAssets(t)

	dnsx_path := writeAssetFixture(t, sandbox, "dnsx", `{"host":"example.com","a":["1.2.3.4"],"cname":["edge.example.net"],"status_code":"NOERROR"}`)
	tlsx_path := writeAssetFixture(t, sandbox, "tlsx", `{"host":"example.com","port":443,"serial":"abc123","subject_cn":"example.com","subject_an":["example.com"],"issuer_cn":"Example CA"}`)

	if _, err := tool.Import(assetArguments{"path": dnsx_path, "format": "dnsx"}); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if _, err := tool.Import(assetArguments{"path": tlsx_path, "format": "tlsx"}); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	entries, err1 := utils_findings.ReadAssetsIndex(tool.Playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	types := map[string]int{}
	parents := 0

	for _, entry := range entries {

		types[string(entry.Type)]++

		if entry.ParentID != "" {
			parents++
		}

	}

	if types["dns_record"] != 1 || types["ip"] != 1 || types["certificate"] != 1 {
		t.Errorf("Expected dns_record, ip and certificate assets, got %v", types)
	}

	if parents != 2 {
		t.Errorf("Expected the IP and certificate assets to have a parent, got %d parent links", parents)
	}

}

func TestAssets_ImportInvalidLines(t *testing.T) {

	tool, sandbox := newTestAssets(t)
	path := writeAssetFixture(t, sandbox, "httpx", "not json\n"+`{"url":"https://example.com/","status_code":200}`)

	result, err0 := tool.Import(assetArguments{"path": path, "format": "httpx"})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err0, result)
	}

	if strings.Contains(result, "skipped 1 invalid") == false {
		t.Errorf("Expected one invalid result, got %q", result)
	}

}
