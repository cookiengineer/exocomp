package ast

import "strings"
import "testing"

func TestReadSymbol_Function(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	result := ReadSymbol(source, "FirstFunction", "func")

	if result == nil {
		t.Fatalf("Expected FirstFunction to be found")
	}

	if result.Name != "FirstFunction" {
		t.Errorf("Expected name %q, got %q", "FirstFunction", result.Name)
	}

	if result.Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", result.Type)
	}

	if strings.Contains(result.Body, "func FirstFunction(current int64, added int64) (string, error)") != true {
		t.Errorf("Expected function signature, got: %q", result.Body)
	}

	if strings.Contains(result.Body, `return "", nil`) != true {
		t.Errorf("Expected function body, got: %q", result.Body)
	}

}

func TestReadSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	data.Name = specification.Name
}
`)

	result := ReadSymbol(source, "Data.Parse", "func")

	if result == nil {
		t.Fatalf("Expected Data.Parse to be found")
	}

	if result.Name != "Data.Parse" {
		t.Errorf("Expected name %q, got %q", "Data.Parse", result.Name)
	}

	if strings.Contains(result.Body, "func (data *Data) Parse(specification *Input)") != true {
		t.Errorf("Expected method signature, got: %q", result.Body)
	}

	if strings.Contains(result.Body, "data.Name = specification.Name") != true {
		t.Errorf("Expected method body, got: %q", result.Body)
	}

}

func TestReadSymbol_Struct(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
	Age  int64
}
`)

	result := ReadSymbol(source, "Data", "struct")

	if result == nil {
		t.Fatalf("Expected Data to be found")
	}

	if result.Name != "Data" {
		t.Errorf("Expected name %q, got %q", "Data", result.Name)
	}

	if strings.Contains(result.Body, "type Data struct") != true {
		t.Errorf("Expected struct declaration, got: %q", result.Body)
	}

	if strings.Contains(result.Body, "Name") != true {
		t.Errorf("Expected struct field, got: %q", result.Body)
	}

	if strings.Contains(result.Body, "Age") != true {
		t.Errorf("Expected struct field, got: %q", result.Body)
	}

}

func TestReadSymbol_Interface(t *testing.T) {

	source := []byte(`package structs

type Parser interface {
	Parse() error
}
`)

	result := ReadSymbol(source, "Parser", "interface")

	if result == nil {
		t.Fatalf("Expected Parser to be found")
	}

	if result.Name != "Parser" {
		t.Errorf("Expected name %q, got %q", "Parser", result.Name)
	}

	if strings.Contains(result.Body, "Parse() error") != true {
		t.Errorf("Expected interface method, got: %q", result.Body)
	}

}

func TestReadSymbol_Missing(t *testing.T) {

	source := []byte(`package core

func FirstFunction() {}
`)

	if result := ReadSymbol(source, "Missing", "func"); result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}

}

func TestReadSymbol_InvalidSource(t *testing.T) {

	if result := ReadSymbol([]byte("func broken"), "broken", "func"); result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}

}
