package ast

import "strings"
import "testing"

func TestSearchSymbols_Function(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	results := SearchSymbols(source, "FirstFunction")

	if len(results) != 1 {
		t.Fatalf("Expected %d results, got %d", 1, len(results))
	}

	if results[0].Name != "FirstFunction" {
		t.Errorf("Expected name %q, got %q", "FirstFunction", results[0].Name)
	}

	if results[0].Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", results[0].Type)
	}

	if strings.Contains(results[0].Body, "func FirstFunction(current int64, added int64) (string, error)") != true {
		t.Errorf("Expected function declaration, got: %q", results[0].Body)
	}

}

func TestSearchSymbols_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	return
}
`)

	results := SearchSymbols(source, "Parse")

	if len(results) != 1 {
		t.Fatalf("Expected %d results, got %d", 1, len(results))
	}

	if results[0].Name != "Data.Parse" {
		t.Errorf("Expected name %q, got %q", "Data.Parse", results[0].Name)
	}

	if results[0].Type != "func" {
		t.Errorf("Expected type %q, got %q", "func", results[0].Type)
	}

	if strings.Contains(results[0].Body, "func (data *Data) Parse(specification *Input)") != true {
		t.Errorf("Expected method declaration, got: %q", results[0].Body)
	}

}

func TestSearchSymbols_Type(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	results := SearchSymbols(source, "Data")

	if len(results) != 1 {
		t.Fatalf("Expected %d results, got %d", 1, len(results))
	}

	if results[0].Name != "Data" {
		t.Errorf("Expected name %q, got %q", "Data", results[0].Name)
	}

	if results[0].Type != "type" {
		t.Errorf("Expected type %q, got %q", "type", results[0].Type)
	}

	if strings.Contains(results[0].Body, "type Data struct") != true {
		t.Errorf("Expected struct declaration, got: %q", results[0].Body)
	}

}

func TestSearchSymbols_CaseInsensitive(t *testing.T) {

	source := []byte(`package core

func ParseInput() {}
`)

	results := SearchSymbols(source, "parseinput")

	if len(results) != 1 {
		t.Fatalf("Expected %d results, got %d", 1, len(results))
	}

	if results[0].Name != "ParseInput" {
		t.Errorf("Expected name %q, got %q", "ParseInput", results[0].Name)
	}

}

func TestSearchSymbols_MultipleMatches(t *testing.T) {

	source := []byte(`package core

func ParseInput() {
	ParseOutput()
}

func ParseOutput() {}
`)

	results := SearchSymbols(source, "ParseOutput")

	if len(results) != 2 {
		t.Fatalf("Expected %d results, got %d", 2, len(results))
	}

}

func TestSearchSymbols_NoMatch(t *testing.T) {

	source := []byte(`package core

func FirstFunction() {}
`)

	results := SearchSymbols(source, "Missing")

	if len(results) != 0 {
		t.Errorf("Expected %d results, got %d", 0, len(results))
	}

}

func TestSearchSymbols_InvalidSource(t *testing.T) {

	results := SearchSymbols([]byte("func broken"), "broken")

	if len(results) != 0 {
		t.Errorf("Expected %d results, got %d", 0, len(results))
	}

}
