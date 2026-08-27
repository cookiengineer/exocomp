package tools

import "encoding/json"
import "os"
import "path/filepath"
import "strings"
import "testing"

import "exocomp/types"

func TestRequirements_DefineFunc(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		result1, err1 := tool.DefineFunc("./core/FirstFunction.go", "FirstFunction", "func FirstFunction(current int64, added int64) (string, error)", "The method needs to implement a fibonacci sequence.")
		result2, err2 := tool.DefineFunc("./parsers/Parse.go", "Parse", "func Parse(specification *structs.Specification, debug bool) *schemas.Result", "The method needs to implement a specification parser.")
		result3, err3 := tool.DefineFunc("./processors/ProcessData.go", "ProcessData", "func ProcessData(specification *structs.Data)", "The method needs to implement a data processor.")
		result4, err4 := tool.DefineFunc("./structs/Data.go", "structs.Data.Parse", "func (data *structs.Data) Parse(specification *schemas.Input)", "The method needs to implement a schema parser.")
		result5, err5 := tool.DefineFunc("./invalid/FunctionWithInvalidType.go", "FunctionWithInvalidType", "function FunctionWithInvalidType (a int, b_custom float1337) (null)", "This method has an invalid b parameter.")

		if result1 != "requirements.DefineFunc: FirstFunction defined as func FirstFunction(current int64, added int64) (string, error)" {
			t.Errorf("Expected FirstFunction to be defined")
		}

		if result2 != "requirements.DefineFunc: Parse defined as func Parse(specification *structs.Specification, debug bool) *schemas.Result" {
			t.Errorf("Expected Parse to be defined")
		}

		if result3 != "requirements.DefineFunc: ProcessData defined as func ProcessData(specification *structs.Data)" {
			t.Errorf("Expected ProcessData to be defined")
		}

		if result4 != "requirements.DefineFunc: structs.Data.Parse defined as func (data *structs.Data) Parse(specification *schemas.Input)" {
			t.Errorf("Expected (*structs.Data) Parse to be defined")
		}

		if result5 != "" {
			t.Errorf("Expected function to be invalid")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if err5 == nil {
			t.Errorf("Expected %v to be not nil", err5)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_DefineInterface(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox := filepath.Join(playground, "requirements")
	tool := NewRequirements(playground, sandbox)

	if tool != nil {

		declaration1 := strings.Join([]string{
			"type Data interface {",
			"\tParse(specification *schemas.Input)",
			"}",
		}, "\n")

		declaration2 := "func (data *structs.Data) Parse(specification *schemas.Input)"

		result1, err1 := tool.DefineInterface("./structs/Data.go", "Data", declaration1, "The interface needs to define a parser.")
		result2, err2 := tool.DefineFunc("./structs/Data.go", "structs.Data.Parse", declaration2, "The method needs to implement a schema parser.")
		result3, err3 := tool.DefineFunc("./structs/Data.go", "DifferentSymbol", declaration1, "The method needs to implement a schema parser.")

		if strings.HasPrefix(result1, "requirements.DefineInterface: Data defined as type Data interface {") == false {
			t.Errorf("Expected \"%s\" to be defined as \"%s\"", "Data", declaration1)
		}

		if result2 != "requirements.DefineFunc: structs.Data.Parse defined as func (data *structs.Data) Parse(specification *schemas.Input)" {
			t.Errorf("Expected \"%s\" to be defined as \"%s\"", "Parse", declaration2)
		}

		if result3 != "" {
			t.Errorf("Expected interface to be invalid")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_DefineStruct(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		declaration1 := strings.Join([]string{
			"type Data struct {",
			"\tName string `json:\"name\"`",
			"\tAge int `json:\"age\"`",
			"\tAddress []string `json:\"address\"`",
			"}",
		}, "\n")

		declaration2 := "func (data *structs.Data) Parse(specification *schemas.Input)"

		result1, err1 := tool.DefineStruct("./structs/Data.go", "Data", declaration1, "The struct needs to implement a database entry for a person.")
		result2, err2 := tool.DefineFunc("./structs/Data.go", "structs.Data.Parse", declaration2, "The method needs to implement a schema parser.")
		result3, err3 := tool.DefineFunc("./structs/Data.go", "DifferentSymbol", declaration1, "The method needs to implement a schema parser.")

		if strings.HasPrefix(result1, "requirements.DefineStruct: Data defined as type Data struct {") == false {
			t.Errorf("Expected \"%s\" to be defined as \"%s\"", "Data", declaration1)
		}

		if result2 != "requirements.DefineFunc: structs.Data.Parse defined as func (data *structs.Data) Parse(specification *schemas.Input)" {
			t.Errorf("Expected \"%s\" to be defined as \"%s\"", "Parse", declaration2)
		}

		if result3 != "" {
			t.Errorf("Expected struct to be invalid")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		declaration1 := strings.Join([]string{
			"type Data struct {",
			"\tName string `json:\"name\"`",
			"\tAge int `json:\"age\"`",
			"\tAddress []string `json:\"address\"`",
			"}",
		}, "\n")

		declaration2 := "func (data *structs.Data) Parse(specification *schemas.Input)"

		_, err1 := tool.DefineFunc("./core/FirstFunction.go", "FirstFunction", "func FirstFunction(current int64, added int64) (string, error)", "The method needs to implement a fibonacci sequence.")
		_, err2 := tool.DefineFunc("./parsers/Parse.go", "Parse", "func Parse(specification *structs.Specification, debug bool) *schemas.Result", "The method needs to implement a specification parser.")
		_, err3 := tool.DefineFunc("./processors/ProcessData.go", "ProcessData", "func ProcessData(specification *structs.Data)", "The method needs to implement a data processor.")
		_, err4 := tool.DefineStruct("./structs/Data.go", "Data", declaration1, "The struct needs to implement a database entry for a person.")
		_, err5 := tool.DefineFunc("./structs/Data.go", "structs.Data.Parse", declaration2, "The method needs to implement a schema parser.")
		result, err6 := tool.List()

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

		if strings.HasPrefix(result, "requirements.List: 5 specifications.") == false{
			t.Errorf("Expected 5 specifications:\n%s", result)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_Search(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		declaration1 := strings.Join([]string{
			"type Data struct {",
			"\tName string `json:\"name\"`",
			"\tAge int `json:\"age\"`",
			"\tAddress []string `json:\"address\"`",
			"}",
		}, "\n")

		declaration2 := "func (data *structs.Data) Parse(specification *schemas.Input)"

		_, err1 := tool.DefineFunc("./core/FirstFunction.go", "FirstFunction", "func FirstFunction(current int64, added int64) (string, error)", "The method needs to implement a fibonacci sequence.")
		_, err2 := tool.DefineFunc("./parsers/Parser.go", "Parse", "func Parse(specification *structs.Specification, debug bool) *schemas.Result", "The method needs to implement a specification parser.")
		_, err3 := tool.DefineFunc("./parsers/Parser.go", "ProcessData", "func ProcessData(specification *structs.Data)", "The method needs to implement a data processor.")
		_, err4 := tool.DefineStruct("./structs/Data.go", "Data", declaration1, "The struct needs to implement a database entry for a person.")
		_, err5 := tool.DefineFunc("./structs/Data.go", "structs.Data.Parse", declaration2, "The method needs to implement a schema parser.")

		result1, err6 := tool.Search("./structs", "")
		result2, err7 := tool.Search("./parsers/Parser.go", "")
		result3, err8 := tool.Search("./parsers/Parser.go", "Process")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if strings.HasPrefix(result1, "requirements.Search: ./structs#* contains 2 specifications.") == false {
			t.Errorf("Expected %d requirement specifications for \"%s\"", 2, "./structs")
		}

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

		if strings.HasPrefix(result2, "requirements.Search: ./parsers/Parser.go#* contains 2 specifications.") == false {
			t.Errorf("Expected %d requirement specifications for \"%s\"", 2, "./parsers/Parser.go")
		}

		if err7 != nil {
			t.Errorf("Expected %v to be nil", err7)
		}

		if strings.HasPrefix(result3, "requirements.Search: ./parsers/Parser.go#Process* contains 1 specifications.") == false {
			t.Errorf("Expected %d requirement specifications for \"%s\"", 1, "./parsers/Parser.go")
		}

		if err8 != nil {
			t.Errorf("Expected %v to be nil", err8)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_Call(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		result1, err1 := tool.Call("DefineFunc", map[string]interface{}{
			"path":        "./structs/Data.go",
			"symbol":      "structs.Data.Parse",
			"declaration": "func (data *structs.Data) Parse(specification *schemas.Input)",
			"behavior":    "The method needs to implement a schema parser.",
		})

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result1 != "requirements.DefineFunc: structs.Data.Parse defined as func (data *structs.Data) Parse(specification *schemas.Input)" {
			t.Errorf("Expected the method to be defined, got %s", result1)
		}

		result2, err2 := tool.Call("DefineStruct", map[string]interface{}{
			"path":        "./structs/Data.go",
			"symbol":      "Data",
			"declaration": "type Data struct {\n\tName string `json:\"name\"`\n}",
			"behavior":    "The struct needs to implement a database entry for a person.",
		})

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if !strings.HasPrefix(result2, "requirements.DefineStruct: Data defined as type Data struct {") {
			t.Errorf("Expected the struct to be defined, got %s", result2)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestRequirements_Signoff(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-requirements-*")
	sandbox       := filepath.Join(playground, "requirements")
	tool          := NewRequirements(playground, sandbox)

	if tool != nil {

		_, err1 := tool.DefineFunc("./core/FirstFunction.go", "FirstFunction", "func FirstFunction(current int64, added int64) (string, error)", "The method needs to implement a fibonacci sequence.")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		source := strings.Join([]string{
			"package core",
			"",
			"func FirstFunction(current int64, added int64) (string, error) {",
			"\treturn \"\", nil",
			"}",
		}, "\n")

		file_path := filepath.Join(sandbox, "core", "FirstFunction.go")
		err2      := os.MkdirAll(filepath.Dir(file_path), 0755)

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		err3 := os.WriteFile(file_path, []byte(source), 0666)

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		result, err4 := tool.Signoff("./core/FirstFunction.go", "FirstFunction")

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if result != "requirements.Signoff: ./core/FirstFunction.go#FirstFunction marked as implemented." {
			t.Errorf("Expected signoff to be successful:\n%s", result)
		}

		reports, err5 := json.Marshal(tool)

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		contents := make(map[string]map[string]types.Requirement)
		err6     := json.Unmarshal(reports, &contents)

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

		reports_count := 0
		reports_implemented := false

		for _, specifications := range contents {
			for _, specification := range specifications {
				reports_count++
				reports_implemented = specification.IsImplemented
			}
		}

		if reports_count != 1 {
			t.Errorf("Expected 1 report, got %d", reports_count)
		} else if reports_implemented != true {
			t.Errorf("Expected report to be implemented")
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
