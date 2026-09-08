package ast

import "testing"

func TestHasSymbol_Function(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	if HasSymbol(source, "FirstFunction", "func") != true {
		t.Errorf("Expected FirstFunction to be found")
	}

	if HasSymbol(source, "MissingFunction", "func") != false {
		t.Errorf("Expected MissingFunction to be not found")
	}

}

func TestHasSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	return
}
`)

	if HasSymbol(source, "Data.Parse", "func") != true {
		t.Errorf("Expected Data.Parse to be found")
	}

	if HasSymbol(source, "Data.Missing", "func") != false {
		t.Errorf("Expected Data.Missing to be not found")
	}

}

func TestHasSymbol_Type(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	if HasSymbol(source, "Data", "struct") != true {
		t.Errorf("Expected Data to be found")
	}

	if HasSymbol(source, "Data", "interface") != false {
		t.Errorf("Expected Data to be not found as interface")
	}

	if HasSymbol(source, "Missing", "struct") != false {
		t.Errorf("Expected Missing to be not found")
	}

}

func TestHasSymbol_BasicTypes(t *testing.T) {

	source := []byte(`package dummy

type MyByte uint8
`)

	if HasSymbol(source, "MyByte", "uint8") != true {
		t.Errorf("Expected MyByte to be found as uint8")
	}

	if HasSymbol(source, "MyByte", "type") != true {
		t.Errorf("Expected MyByte to be found as type")
	}

	if HasSymbol(source, "MyByte", "") != true {
		t.Errorf("Expected MyByte to be found with empty type")
	}

	if HasSymbol(source, "Missing", "uint8") != false {
		t.Errorf("Expected Missing to be not found")
	}

}
