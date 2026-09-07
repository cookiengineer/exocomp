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

	if strings.Contains(result, "func FirstFunction(current int64, added int64) (string, error)") != true {
		t.Errorf("Expected function signature, got: %q", result)
	}

	if strings.Contains(result, `return "", nil`) != true {
		t.Errorf("Expected function body, got: %q", result)
	}

}

func TestReadSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	data.Name = specification.Name
}
`)

	result := ReadSymbol(source, "Data.Parse", "func")

	if strings.Contains(result, "func (data *Data) Parse(specification *Input)") != true {
		t.Errorf("Expected method signature, got: %q", result)
	}

	if strings.Contains(result, "data.Name = specification.Name") != true {
		t.Errorf("Expected method body, got: %q", result)
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

	if strings.Contains(result, "type Data struct") != true {
		t.Errorf("Expected struct declaration, got: %q", result)
	}

	if strings.Contains(result, "Name") != true {
		t.Errorf("Expected struct field, got: %q", result)
	}

	if strings.Contains(result, "Age") != true {
		t.Errorf("Expected struct field, got: %q", result)
	}

}

func TestReadSymbol_Interface(t *testing.T) {

	source := []byte(`package structs

type Parser interface {
	Parse() error
}
`)

	result := ReadSymbol(source, "Parser", "interface")

	if strings.Contains(result, "Parse() error") != true {
		t.Errorf("Expected interface method, got: %q", result)
	}

}

func TestReadSymbol_Missing(t *testing.T) {

	source := []byte(`package core

func FirstFunction() {}
`)

	if result := ReadSymbol(source, "Missing", "func"); result != "" {
		t.Errorf("Expected empty result, got: %q", result)
	}

}

func TestReadSymbol_InvalidSource(t *testing.T) {

	if result := ReadSymbol([]byte("func broken"), "broken", "func"); result != "" {
		t.Errorf("Expected empty result, got: %q", result)
	}

}
