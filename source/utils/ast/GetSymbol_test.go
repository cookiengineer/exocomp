package ast

import "strings"
import "testing"

func TestGetSymbol_Function(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	result := GetSymbol(source, "FirstFunction", "func")

	if result != "func FirstFunction(current int64, added int64) (string, error)" {
		t.Errorf("Expected function header, got: %q", result)
	}

}

func TestGetSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	return
}
`)

	result := GetSymbol(source, "Data.Parse", "func")

	if result != "func (data *Data) Parse(specification *Input)" {
		t.Errorf("Expected method header, got: %q", result)
	}

}

func TestGetSymbol_Struct(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	result := GetSymbol(source, "Data", "struct")

	if result != "type Data struct" {
		t.Errorf("Expected struct header, got: %q", result)
	}

}

func TestGetSymbol_Interface(t *testing.T) {

	source := []byte(`package structs

type Parser interface {
	Parse() error
}
`)

	result := GetSymbol(source, "Parser", "interface")

	if result != "type Parser interface" {
		t.Errorf("Expected interface header, got: %q", result)
	}

}

func TestGetSymbol_FuncTypeAlias(t *testing.T) {

	source := []byte(`package core

type Handler func(string) error
`)

	result := GetSymbol(source, "Handler", "func")

	if result != "type Handler func(string) error" {
		t.Errorf("Expected func type alias, got: %q", result)
	}

}

func TestGetSymbol_Missing(t *testing.T) {

	source := []byte(`package core

func FirstFunction() {}
`)

	if result := GetSymbol(source, "Missing", "func"); result != "" {
		t.Errorf("Expected empty result, got: %q", result)
	}

	if result := GetSymbol(source, "FirstFunction", "struct"); result != "" {
		t.Errorf("Expected empty result for wrong type, got: %q", result)
	}

}

func TestGetSymbol_InvalidSource(t *testing.T) {

	if result := GetSymbol([]byte("func broken"), "broken", "func"); result != "" {
		t.Errorf("Expected empty result, got: %q", result)
	}

}

func TestGetSymbol_PrefersFunctionOverType(t *testing.T) {

	source := []byte(`package core

type Handler func(string) error

func Handler() {}
`)

	result := GetSymbol(source, "Handler", "func")

	if strings.Contains(result, "func Handler()") != true {
		t.Errorf("Expected function to be preferred, got: %q", result)
	}

}
