package files

import "os"
import "path/filepath"
import "strings"
import "testing"

func newReservedFilesTool(t *testing.T) (*Files, string) {

	t.Helper()

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "sandbox")

	if err := os.MkdirAll(sandbox, 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.MkdirAll(filepath.Join(sandbox, ".exocomp", "findings"), 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(filepath.Join(sandbox, ".exocomp", "findings", "index.json"), []byte("[]"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(filepath.Join(sandbox, ".exocomp", "hidden.go"), []byte("package hidden\n\nfunc Hidden() {}\n"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	if err := os.WriteFile(filepath.Join(sandbox, "source.txt"), []byte("source"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	tool := NewFiles([]string{"Copy", "List", "Read", "ReadSymbol", "Search", "Stat", "Write", "WriteSymbol"}, playground, sandbox)

	return tool, sandbox

}

func assertReservedError(t *testing.T, method string, result string, err error, sandbox string) {

	t.Helper()

	if err == nil {
		t.Fatalf("Expected %s to reject .exocomp, got %q", method, result)
	}

	if strings.Contains(err.Error(), "reserved") == false {
		t.Errorf("Expected %s to report a reserved path, got %q", method, err.Error())
	}

	if strings.Contains(err.Error(), sandbox) == true {
		t.Errorf("Expected %s error to hide the sandbox path, got %q", method, err.Error())
	}

}

func TestFiles_RejectsExocompFolder(t *testing.T) {

	tool, sandbox := newReservedFilesTool(t)

	result1, err1 := tool.List(".exocomp")
	assertReservedError(t, "files.List", result1, err1, sandbox)

	result2, err2 := tool.Read(".exocomp/findings/index.json")
	assertReservedError(t, "files.Read", result2, err2, sandbox)

	result3, err3 := tool.ReadSymbol(".exocomp/hidden.go", "Hidden")
	assertReservedError(t, "files.ReadSymbol", result3, err3, sandbox)

	result4, err4 := tool.Search(".exocomp", "Hidden")
	assertReservedError(t, "files.Search", result4, err4, sandbox)

	result5, err5 := tool.Stat(".exocomp")
	assertReservedError(t, "files.Stat", result5, err5, sandbox)

	result6, err6 := tool.Write(".exocomp/created.txt", "created")
	assertReservedError(t, "files.Write", result6, err6, sandbox)

	if _, err := os.Stat(filepath.Join(sandbox, ".exocomp", "created.txt")); os.IsNotExist(err) == false {
		t.Errorf("Expected rejected files.Write to create no file")
	}

	result7, err7 := tool.WriteSymbol(".exocomp/hidden.go", "Hidden", "func Hidden() { panic(\"changed\") }")
	assertReservedError(t, "files.WriteSymbol", result7, err7, sandbox)

	result8, err8 := tool.Copy(".exocomp/hidden.go", "./copied.go")
	assertReservedError(t, "files.Copy", result8, err8, sandbox)

	result9, err9 := tool.Copy("./source.txt", ".exocomp/copied.txt")
	assertReservedError(t, "files.Copy", result9, err9, sandbox)

	if _, err := os.Stat(filepath.Join(sandbox, ".exocomp", "copied.txt")); os.IsNotExist(err) == false {
		t.Errorf("Expected rejected files.Copy to create no file")
	}

}

func TestFiles_RejectsNestedAndMixedCaseExocomp(t *testing.T) {

	tool, sandbox := newReservedFilesTool(t)

	cases := []string{
		"./.exocomp/created.txt",
		"./nested/.exocomp/created.txt",
		"./.exocomp/../.exocomp/created.txt",
		"./.EXOCOMP/created.txt",
	}

	for _, path := range cases {

		result, err := tool.Write(path, "created")
		assertReservedError(t, "files.Write", result, err, sandbox)

	}

	if _, err := os.Stat(filepath.Join(sandbox, "nested", ".exocomp", "created.txt")); os.IsNotExist(err) == false {
		t.Errorf("Expected rejected nested files.Write to create no file")
	}

	if _, err := os.Stat(filepath.Join(sandbox, ".EXOCOMP", "created.txt")); os.IsNotExist(err) == false {
		t.Errorf("Expected rejected mixed-case files.Write to create no file")
	}

}

func TestFiles_ExocompIsFilteredFromListAndSearch(t *testing.T) {

	tool, _ := newReservedFilesTool(t)

	list, err0 := tool.List(".")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	if strings.Contains(list, ".exocomp") == true {
		t.Errorf("Expected .exocomp to be filtered from files.List, got %q", list)
	}

	search, err1 := tool.Search(".", "Hidden")

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if strings.Contains(search, ".exocomp") == true {
		t.Errorf("Expected .exocomp to be filtered from files.Search, got %q", search)
	}

}
