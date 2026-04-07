package tools

import "exocomp/utils"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "bytes"
import "fmt"
import "strings"

type Requirements struct {
	Sandbox string
}

func NewRequirements(agent string, sandbox string) *Requirements {

	// TODO: Requirements need to be specified in different folder

	requirements := &Requirements{
		Sandbox: sandbox,
	}

	return requirements

}

func (tool *Requirements) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "DefineFunc" {

		// TODO
		return "", nil

	} else if method == "DefineStruct" {

		// TODO
		return "", nil

	} else if method == "DefineTest" {

		// TODO
		return "", nil

	} else if method == "Search" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)

		if ok1 == true && ok2 == true {
			return tool.Search(utils.FormatFilePath(path), utils.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils.FormatFilePath(path), "")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("requirements.%s: Invalid method.", method)
	}

}

func (tool *Requirements) List() (string, error) {
	return "", nil
}

func (tool *Requirements) DefineFunc(path string, symbol string, definition string, description string) (string, error) {

	fileset    := token.NewFileSet()
	file, err0 := parser.ParseFile(fileset, "", strings.Join([]string{
		"package dummy",
		definition,
	}, "\n"), 0)

	if err0 == nil {

		var found *ast.FuncDecl = nil

		for _, decl := range file.Decls {

			declaration, ok0 := decl.(*ast.FuncDecl)

			if ok0 == true {

				if declaration.Type.Params != nil && len(declaration.Type.Params.List) > 0 {

					buffer := bytes.Buffer{}
					printer.Fprint(&buffer, token.NewFileSet(), decl)

					fmt.Println(buffer.String())

				}

			}

		}

		if found != nil {

			// TODO: Write to a file

		} else {
		}

	}

	return "", nil

}

func (tool *Requirements) DefineStruct(path string, symbol string, fields []string, description string) (string, error) {
	return "", nil
}

func (tool *Requirements) DefineTest(path string, symbol string, inputs []string, outputs []string, description string) (string, error) {
	return "", nil
}

func (tool *Requirements) Search(path string, symbol string) (string, error) {
	return "", nil
}
