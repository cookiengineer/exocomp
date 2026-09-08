package ast

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestGetPackageSymbols(t *testing.T) {

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

	result := GetPackageSymbols(folder, false)

	key1 := file1
	key2 := file2

	symbols1, ok1 := result[key1]

	if ok1 != true {
		t.Errorf("Expected %q to be present in result, got keys: %v", key1, result)
	}

	symbols2, ok2 := result[key2]

	if ok2 != true {
		t.Errorf("Expected %q to be present in result, got keys: %v", key2, result)
	}

	first, ok3 := symbols1["FirstFunction"]

	if ok3 != true {
		t.Fatalf("Expected FirstFunction to be present in %q", key1)
	}

	if first.Name != "FirstFunction" {
		t.Errorf("Expected name %q, got %q", "FirstFunction", first.Name)
	}

	if first.Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", first.Type)
	}

	if strings.Contains(first.Body, "func FirstFunction()") != true {
		t.Errorf("Expected FirstFunction header, got: %q", first.Body)
	}

	data, ok4 := symbols1["Data"]

	if ok4 != true {
		t.Fatalf("Expected Data to be present in %q", key1)
	}

	if data.Body != "type Data struct" {
		t.Errorf("Expected Data header, got: %q", data.Body)
	}

	parse, ok5 := symbols2["Data.Parse"]

	if ok5 != true {
		t.Fatalf("Expected Data.Parse to be present in %q", key2)
	}

	if parse.Name != "Data.Parse" {
		t.Errorf("Expected name %q, got %q", "Data.Parse", parse.Name)
	}

	if parse.Body != "func (data *Data) Parse() error" {
		t.Errorf("Expected Data.Parse header, got: %q", parse.Body)
	}

}

func TestGetPackageSymbols_BasicTypes(t *testing.T) {

	folder := t.TempDir()

	file1 := filepath.Join(folder, "types.go")

	os.WriteFile(file1, []byte(`package core

type MyByte uint8

type MyString string

type MyBytes []byte
`), 0644)

	result := GetPackageSymbols(folder, false)

	key1 := file1

	symbols1, ok1 := result[key1]

	if ok1 != true {
		t.Fatalf("Expected %q to be present in result, got keys: %v", key1, result)
	}

	my_byte, ok2 := symbols1["MyByte"]

	if ok2 != true {
		t.Fatalf("Expected MyByte to be present in %q", key1)
	}

	if my_byte.Name != "MyByte" {
		t.Errorf("Expected name %q, got %q", "MyByte", my_byte.Name)
	}

	if my_byte.Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", my_byte.Type)
	}

	if my_byte.Body != "type MyByte uint8" {
		t.Errorf("Expected body %q, got %q", "type MyByte uint8", my_byte.Body)
	}

	my_string, ok3 := symbols1["MyString"]

	if ok3 != true {
		t.Fatalf("Expected MyString to be present in %q", key1)
	}

	if my_string.Body != "type MyString string" {
		t.Errorf("Expected body %q, got %q", "type MyString string", my_string.Body)
	}

	my_bytes, ok4 := symbols1["MyBytes"]

	if ok4 != true {
		t.Fatalf("Expected MyBytes to be present in %q", key1)
	}

	if my_bytes.Body != "type MyBytes []byte" {
		t.Errorf("Expected body %q, got %q", "type MyBytes []byte", my_bytes.Body)
	}

}

func TestGetPackageSymbols_EmptyFolder(t *testing.T) {

	folder := t.TempDir()

	if result := GetPackageSymbols(folder, false); len(result) != 0 {
		t.Errorf("Expected empty result, got: %v", result)
	}

}

func TestGetPackageSymbols_IgnoresNonGoFiles(t *testing.T) {

	folder := t.TempDir()

	os.WriteFile(filepath.Join(folder, "README.md"), []byte("# readme"), 0644)

	if result := GetPackageSymbols(folder, false); len(result) != 0 {
		t.Errorf("Expected empty result for non-go files, got: %v", result)
	}

}

func TestGetPackageSymbols_IgnoresInvalidGo(t *testing.T) {

	folder := t.TempDir()

	os.WriteFile(filepath.Join(folder, "broken.go"), []byte("this is not go"), 0644)
	os.WriteFile(filepath.Join(folder, "valid.go"), []byte("package core\n\nfunc Valid() {}\n"), 0644)

	result := GetPackageSymbols(folder, false)

	if len(result) != 1 {
		t.Errorf("Expected only valid go file, got: %v", result)
	}

}
