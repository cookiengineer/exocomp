package ast

import "strings"
import "testing"

func TestWriteSymbol_Function(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	declaration := `func FirstFunction(current int64, added int64) (string, error) {
	return "changed", nil
}`

	result := WriteSymbol(source, "FirstFunction", declaration, "func")
	text := string(result)

	if strings.Contains(text, `return "changed", nil`) != true {
		t.Errorf("Expected FirstFunction body to be overwritten, got: %s", text)
	}

	if strings.Contains(text, `return "", nil`) == true {
		t.Errorf("Expected old FirstFunction body to be gone, got: %s", text)
	}

}

func TestWriteSymbol_Method(t *testing.T) {

	source := []byte(`package structs

func (data *Data) Parse(specification *Input) {
	return
}
`)

	declaration := `func (data *Data) Parse(specification *Input) {
	data.Name = specification.Name
}`

	result := WriteSymbol(source, "Data.Parse", declaration, "func")
	text := string(result)

	if strings.Contains(text, "data.Name = specification.Name") != true {
		t.Errorf("Expected Data.Parse body to be overwritten, got: %s", text)
	}

	if strings.Contains(text, "return") == true {
		t.Errorf("Expected old Data.Parse body to be gone, got: %s", text)
	}

}

func TestWriteSymbol_Struct(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	declaration := `type Data struct {
	Name string
	Age  int64
}`

	result := WriteSymbol(source, "Data", declaration, "struct")
	text := string(result)

	if strings.Contains(text, "Age") != true {
		t.Errorf("Expected Data struct to be overwritten, got: %s", text)
	}

	if strings.Contains(text, "Name") != true {
		t.Errorf("Expected Data struct to keep Name field, got: %s", text)
	}

}

func TestWriteSymbol_Interface(t *testing.T) {

	source := []byte(`package structs

type Parser interface {
	Parse() error
}
`)

	declaration := `type Parser interface {
	Parse() error
	Close() error
}`

	result := WriteSymbol(source, "Parser", declaration, "interface")
	text := string(result)

	if strings.Contains(text, "Close() error") != true {
		t.Errorf("Expected Parser interface to be overwritten, got: %s", text)
	}

}

func TestWriteSymbol_BasicTypes(t *testing.T) {

	source := []byte(`package dummy

type MyByte uint8
`)

	declaration := `type MyByte string`

	result := WriteSymbol(source, "MyByte", declaration, "uint8")
	text := string(result)

	if strings.Contains(text, "type MyByte string") != true {
		t.Errorf("Expected MyByte type to be overwritten, got: %s", text)
	}

	if strings.Contains(text, "uint8") == true {
		t.Errorf("Expected old MyByte type to be gone, got: %s", text)
	}

}

func TestWriteSymbol_MissingSymbol(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	declaration := `func MissingFunction(current int64) (string, error) {
	return "changed", nil
}`

	result := WriteSymbol(source, "MissingFunction", declaration, "func")

	if string(result) != string(source) {
		t.Errorf("Expected source to be unchanged for missing symbol, got: %s", string(result))
	}

}

func TestWriteSymbol_NameMismatch(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	declaration := `func OtherFunction(current int64) (string, error) {
	return "changed", nil
}`

	result := WriteSymbol(source, "FirstFunction", declaration, "func")

	if string(result) != string(source) {
		t.Errorf("Expected source to be unchanged on name mismatch, got: %s", string(result))
	}

}

func TestWriteSymbol_InvalidDeclaration(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}
`)

	result := WriteSymbol(source, "FirstFunction", "func FirstFunction(", "func")

	if string(result) != string(source) {
		t.Errorf("Expected source to be unchanged on invalid declaration, got: %s", string(result))
	}

}

func TestWriteSymbol_TypeMismatch(t *testing.T) {

	source := []byte(`package structs

type Data struct {
	Name string
}
`)

	declaration := `type Data interface {
	Parse() error
}`

	result := WriteSymbol(source, "Data", declaration, "interface")

	if string(result) != string(source) {
		t.Errorf("Expected source to be unchanged on type mismatch, got: %s", string(result))
	}

}

func TestWriteSymbol_FuncTypeAlias(t *testing.T) {

	source := []byte(`package core

type Handler func(string) error
`)

	declaration := `type Handler func(string, int64) error`

	result := WriteSymbol(source, "Handler", declaration, "func")
	text := string(result)

	if strings.Contains(text, "func(string, int64) error") != true {
		t.Errorf("Expected Handler func type to be overwritten, got: %s", text)
	}

	if strings.Contains(text, "func(string) error") == true {
		t.Errorf("Expected old Handler func type to be gone, got: %s", text)
	}

}

func TestWriteSymbol_PreservesOtherSymbols(t *testing.T) {

	source := []byte(`package core

func FirstFunction(current int64, added int64) (string, error) {
	return "", nil
}

func SecondFunction() {}

type Data struct {
	Name string
}
`)

	declaration := `func FirstFunction(current int64, added int64) (string, error) {
	return "changed", nil
}`

	result := WriteSymbol(source, "FirstFunction", declaration, "func")
	text := string(result)

	if strings.Contains(text, "SecondFunction") != true {
		t.Errorf("Expected SecondFunction to be preserved, got: %s", text)
	}

	if strings.Contains(text, "type Data struct") != true {
		t.Errorf("Expected Data struct to be preserved, got: %s", text)
	}

}
