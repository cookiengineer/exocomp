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
		t.Errorf("Expected function declaration, got: %q", result.Body)
	}

}

func TestGetSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	return
}
`)

	result := GetSymbol(source, "Data.Parse", "func")

	if result == nil {
		t.Fatalf("Expected Data.Parse to be found")
	}

	if result.Name != "Data.Parse" {
		t.Errorf("Expected name %q, got %q", "Data.Parse", result.Name)
	}

	if result.Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", result.Type)
	}

	if strings.Contains(result.Body, "func (data *Data) Parse(specification *Input)") != true {
		t.Errorf("Expected method declaration, got: %q", result.Body)
	}

}

func TestGetSymbol_Struct(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	result := GetSymbol(source, "Data", "struct")

	if result == nil {
		t.Fatalf("Expected Data to be found")
	}

	if result.Name != "Data" {
		t.Errorf("Expected name %q, got %q", "Data", result.Name)
	}

	if result.Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", result.Type)
	}

	if strings.Contains(result.Body, "type Data struct") != true {
		t.Errorf("Expected struct declaration, got: %q", result.Body)
	}

}

func TestGetSymbol_Interface(t *testing.T) {

	source := []byte(`package structs

type Parser interface {
	Parse() error
}
`)

	result := GetSymbol(source, "Parser", "interface")

	if result == nil {
		t.Fatalf("Expected Parser to be found")
	}

	if result.Name != "Parser" {
		t.Errorf("Expected name %q, got %q", "Parser", result.Name)
	}

	if result.Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", result.Type)
	}

	if strings.Contains(result.Body, "type Parser interface") != true {
		t.Errorf("Expected interface declaration, got: %q", result.Body)
	}

}

func TestGetSymbol_BasicTypes(t *testing.T) {

	tests := []struct {
		source          string
		symbol          string
		declaration_type string
		want            string
	}{
		{"package dummy\ntype MyByte uint8", "MyByte", "uint8", "type MyByte uint8"},
		{"package dummy\ntype MyString string", "MyString", "string", "type MyString string"},
		{"package dummy\ntype MyBytes []byte", "MyBytes", "[]byte", "type MyBytes []byte"},
		{"package dummy\ntype IPv4 [4]byte", "IPv4", "[4]byte", "type IPv4 [4]byte"},
		{"package dummy\ntype Flags []bool", "Flags", "[]bool", "type Flags []bool"},
		{"package dummy\ntype Counts map[string]int", "Counts", "map[string]int", "type Counts map[string]int"},
		{"package dummy\ntype ScaleFactor float64", "ScaleFactor", "float64", "type ScaleFactor float64"},
	}

	for _, test := range tests {

		result := GetSymbol([]byte(test.source), test.symbol, test.declaration_type)

		if result == nil {
			t.Fatalf("Expected %q to be found", test.symbol)
		}

		if result.Name != test.symbol {
			t.Errorf("Expected name %q, got %q", test.symbol, result.Name)
		}

		if result.Type != "type" {
			t.Errorf("Expected type %q, got %q", "type", result.Type)
		}

		if strings.Contains(result.Body, test.want) != true {
			t.Errorf("Expected body to contain %q, got %q", test.want, result.Body)
		}

	}

}

func TestGetSymbol_BasicTypes_EmptyDeclarationType(t *testing.T) {

	source := []byte(`package dummy

type MyByte uint8
`)

	result := GetSymbol(source, "MyByte", "")

	if result == nil {
		t.Fatalf("Expected MyByte to be found")
	}

	if result.Name != "MyByte" {
		t.Errorf("Expected name %q, got %q", "MyByte", result.Name)
	}

	if result.Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", result.Type)
	}

	if strings.Contains(result.Body, "type MyByte uint8") != true {
		t.Errorf("Expected body to contain %q, got %q", "type MyByte uint8", result.Body)
	}

}

func TestGetSymbol_FuncTypeAlias(t *testing.T) {

	source := []byte(`package core

type Handler func(string) error
`)

	result := GetSymbol(source, "Handler", "func")

	if result == nil {
		t.Fatalf("Expected Handler to be found")
	}

	if result.Name != "Handler" {
		t.Errorf("Expected name %q, got %q", "Handler", result.Name)
	}

	if result.Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", result.Type)
	}

	if strings.Contains(result.Body, "type Handler func(string) error") != true {
		t.Errorf("Expected func type alias, got: %q", result.Body)
	}

}

func TestGetSymbol_Missing(t *testing.T) {

	source := []byte(`package core

func FirstFunction() {}
`)

	if result := GetSymbol(source, "Missing", "func"); result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}

	if result := GetSymbol(source, "FirstFunction", "struct"); result != nil {
		t.Errorf("Expected nil result for wrong type, got: %v", result)
	}

}

func TestGetSymbol_InvalidSource(t *testing.T) {

	if result := GetSymbol([]byte("func broken"), "broken", "func"); result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}

}

func TestGetSymbol_PrefersFunctionOverType(t *testing.T) {

	source := []byte(`package core

type Handler func(string) error

func Handler() {}
`)

	result := GetSymbol(source, "Handler", "func")

	if result == nil {
		t.Fatalf("Expected Handler to be found")
	}

	if result.Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", result.Type)
	}

	if strings.Contains(result.Body, "func Handler()") != true {
		t.Errorf("Expected function to be preferred, got: %q", result.Body)
	}

}
