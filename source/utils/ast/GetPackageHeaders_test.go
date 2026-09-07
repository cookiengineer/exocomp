package ast

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestGetPackageHeaders(t *testing.T) {

	folder := t.TempDir()

	file1 := filepath.Join(folder, "functions.go")

	os.WriteFile(file1, []byte(`package core

func FirstFunction() {}

type Data struct {
	Name string
}
`), 0644)

	file2 := filepath.Join(folder, "methods.go")

	os.WriteFile(file2, []byte(`package core

func (data *Data) Parse() error {
	return nil
}
`), 0644)

	result := GetPackageHeaders(folder)

	key1 := file1
	key2 := file2

	_, ok1 := result[key1]

	if ok1 != true {
		t.Errorf("Expected %q to be present in result, got keys: %v", key1, result)
	}

	_, ok2 := result[key2]

	if ok2 != true {
		t.Errorf("Expected %q to be present in result, got keys: %v", key2, result)
	}

	headers := result[key1]

	if strings.Contains(headers["FirstFunction"], "func FirstFunction()") != true {
		t.Errorf("Expected FirstFunction header, got: %q", headers["FirstFunction"])
	}

	if headers["Data"] != "type Data struct" {
		t.Errorf("Expected Data header, got: %q", headers["Data"])
	}

	headers2 := result[key2]

	if headers2["Data.Parse"] != "func (data *Data) Parse() error" {
		t.Errorf("Expected Data.Parse header, got: %q", headers2["Data.Parse"])
	}

}

func TestGetPackageHeaders_EmptyFolder(t *testing.T) {

	folder := t.TempDir()

	if result := GetPackageHeaders(folder); len(result) != 0 {
		t.Errorf("Expected empty result, got: %v", result)
	}

}

func TestGetPackageHeaders_IgnoresNonGoFiles(t *testing.T) {

	folder := t.TempDir()

	os.WriteFile(filepath.Join(folder, "README.md"), []byte("# readme"), 0644)

	if result := GetPackageHeaders(folder); len(result) != 0 {
		t.Errorf("Expected empty result for non-go files, got: %v", result)
	}

}

func TestGetPackageHeaders_IgnoresInvalidGo(t *testing.T) {

	folder := t.TempDir()

	os.WriteFile(filepath.Join(folder, "broken.go"), []byte("this is not go"), 0644)
	os.WriteFile(filepath.Join(folder, "valid.go"), []byte("package core\n\nfunc Valid() {}\n"), 0644)

	result := GetPackageHeaders(folder)

	if len(result) != 1 {
		t.Errorf("Expected only valid go file, got: %v", result)
	}

}
