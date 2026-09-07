package ast

import "testing"

func TestGetSymbolType_BasicTypes(t *testing.T) {

	tests := []struct {
		source string
		symbol string
		want   string
	}{
		{"package dummy\ntype MyEnum string", "MyEnum", "string"},
		{"package dummy\ntype IPv4 [4]byte", "IPv4", "[4]byte"},
		{"package dummy\ntype ScaleFactor float64", "ScaleFactor", "float64"},
		{"package dummy\ntype Flags []bool", "Flags", "[]bool"},
		{"package dummy\ntype Counts map[string]int", "Counts", "map[string]int"},
	}

	for _, test := range tests {

		result, err := GetSymbolType([]byte(test.source), test.symbol)

		if err != nil {
			t.Errorf("Expected %v to be nil, got %v", nil, err)
		}

		if result != test.want {
			t.Errorf("Expected %q, got %q", test.want, result)
		}

	}

}

func TestGetSymbolType_RejectsComplexTypes(t *testing.T) {

	tests := []struct {
		source string
		symbol string
	}{
		{"package dummy\ntype Data struct {\n\tName string\n}", "Data"},
		{"package dummy\ntype Parser interface {\n\tParse() error\n}", "Parser"},
		{"package dummy\ntype Handler func(string) error", "Handler"},
	}

	for _, test := range tests {

		result, err := GetSymbolType([]byte(test.source), test.symbol)

		if result != "" {
			t.Errorf("Expected empty result, got %q", result)
		}

		if err == nil {
			t.Errorf("Expected %v to be not nil", nil)
		}

	}

}

func TestGetSymbolType_MissingSymbol(t *testing.T) {

	result, err := GetSymbolType([]byte("package dummy\ntype MyEnum string"), "Missing")

	if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

	if err == nil {
		t.Errorf("Expected %v to be not nil", nil)
	}

}

func TestGetSymbolType_FuncReturnType(t *testing.T) {

	tests := []struct {
		source string
		symbol string
		want   string
	}{
		{"package dummy\nfunc A() (string, error) { return \"\", nil }", "A", "(string, error)"},
		{"package dummy\nfunc B() string { return \"\" }", "B", "string"},
		{"package dummy\nfunc C() int8 { return 0 }", "C", "int8"},
		{"package dummy\nfunc E() (a, b int) { return 0, 0 }", "E", "(int, int)"},
		{"package dummy\nfunc F() (x int, y error) { return 0, nil }", "F", "(int, error)"},
	}

	for _, test := range tests {

		result, err := GetSymbolType([]byte(test.source), test.symbol)

		if err != nil {
			t.Errorf("Expected %v to be nil, got %v", nil, err)
		}

		if result != test.want {
			t.Errorf("Expected %q, got %q", test.want, result)
		}

	}

}

func TestGetSymbolType_FuncMethodReturnType(t *testing.T) {

	result, err := GetSymbolType([]byte("package dummy\nfunc (d *Data) Parse() (string, error) { return \"\", nil }"), "Data.Parse")

	if err != nil {
		t.Errorf("Expected %v to be nil, got %v", nil, err)
	}

	if result != "(string, error)" {
		t.Errorf("Expected %q, got %q", "(string, error)", result)
	}

}

func TestGetSymbolType_FuncWithoutReturnType(t *testing.T) {

	result, err := GetSymbolType([]byte("package dummy\nfunc D() { }"), "D")

	if err != nil {
		t.Errorf("Expected %v to be nil, got %v", nil, err)
	}

	if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

}

func TestGetSymbolType_InvalidSyntax(t *testing.T) {

	result, err := GetSymbolType([]byte("type broken"), "broken")

	if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

	if err == nil {
		t.Errorf("Expected %v to be not nil", nil)
	}

}
