package tools

import utils_fmt "exocomp/utils/fmt"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "bytes"
import "fmt"
import "sort"
import "strings"

type requirement_specification struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Symbol      string `json:"symbol"`
	Declaration string `json:"declaration"`
	Behavior    string `json:"behavior"`
}

type Requirements struct {
	Playground string
	Sandbox    string
	contents   map[string]map[string]requirement_specification // map[resolved][symbol]
}

func NewRequirements(playground string, sandbox string) *Requirements {

	requirements := &Requirements{
		Playground: playground,
		Sandbox:    sandbox,
		contents:   make(map[string]map[string]requirement_specification),
	}

	readRequirements(requirements)

	return requirements

}

func (tool *Requirements) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "DefineFunc" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		declaration, ok3 := arguments["declaration"].(string)
		behavior,    ok4 := arguments["behavior"].(string)

		if ok1 == true && ok2 == true && ok3 == true && ok4 == true {
			return tool.DefineFunc(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(declaration), utils_fmt.FormatSingleLine(behavior))
		} else if ok1 == true && ok2 == true && ok3 == true && ok4 == false {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"behavior\" is not a string.")
		} else if ok1 == true && ok2 == true && ok3 == false && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"declaration\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "DefineStruct" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		declaration, ok3 := arguments["declaration"].(string)
		behavior,    ok4 := arguments["behavior"].(string)

		if ok1 == true && ok2 == true && ok3 == true && ok4 == true {
			return tool.DefineFunc(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(declaration), utils_fmt.FormatSingleLine(behavior))
		} else if ok1 == true && ok2 == true && ok3 == true && ok4 == false {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"behavior\" is not a string.")
		} else if ok1 == true && ok2 == true && ok3 == false && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"declaration\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true && ok4 == true {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Search" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)

		if ok1 == true && ok2 == true {
			return tool.Search(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils_fmt.FormatFilePath(path), "")
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

	lines := make([]string, 0)

	for _, specifications := range tool.contents {

		for _, specification := range specifications {

			resolved_path, err1 := resolveSandboxPath(tool.Playground, specification.File)

			if err1 == nil {

				sandbox_path, err2 := sanitizeSandboxPath(tool.Sandbox, resolved_path)

				if err2 == nil {
					lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Declaration: %s, Behavior: %s", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior))
				}

			}

		}

	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("requirements.List: %d specifications.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Requirements) DefineFunc(path string, symbol string, declaration string, behavior string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		fileset    := token.NewFileSet()
		file, err1 := parser.ParseFile(fileset, "", strings.Join([]string{
			"package dummy",
			declaration,
		}, "\n"), 0)

		if err1 == nil {

			declaration_symbol := ""
			declaration_code   := ""

			for _, decl := range file.Decls {

				declaration, ok0 := decl.(*ast.FuncDecl)

				if ok0 == true {

					if declaration.Name != nil {

						if declaration.Recv != nil && len(declaration.Recv.List) > 0 {

							recv_type := declaration.Recv.List[0].Type
							type_name := ""

							switch typ := recv_type.(type) {

							case *ast.Ident:
								type_name = typ.Name

							case *ast.StarExpr:

								ident, ok := typ.X.(*ast.Ident)

								if ok == true {
									type_name = ident.Name
								}

							}

							if type_name != "" {
								declaration_symbol = type_name + "." + declaration.Name.Name
							} else {
								declaration_symbol = declaration.Name.Name
							}

						} else {
							declaration_symbol = declaration.Name.Name
						}

					}

					buffer := bytes.Buffer{}
					printer.Fprint(&buffer, token.NewFileSet(), declaration)
					declaration_code = strings.TrimSpace(buffer.String())

					break

				}

			}

			if declaration_symbol == symbol {

				_, ok1 := tool.contents[resolved]

				if ok1 == false {
					tool.contents[resolved] = make(map[string]requirement_specification)
				}

				tool.contents[resolved][symbol] = requirement_specification{
					File:        resolved,
					Type:        "func",
					Declaration: declaration_code,
					Symbol:      declaration_symbol,
					Behavior:    behavior,
				}

				err2 := writeRequirements(tool)

				if err2 == nil {
					return fmt.Sprintf("requirements.DefineFunc: %s defined as %s", declaration_symbol, declaration_code), nil
				} else {
					return "", fmt.Errorf("requirements.DefineFunc: %s", err2.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineFunc: Invalid syntax, function symbol \"%s\" must be the same as symbol \"%s\".", declaration_symbol, symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineFunc: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineFunc: %s", err0.Error())
	}

}

func (tool *Requirements) DefineStruct(path string, symbol string, declaration string, behavior string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		fileset    := token.NewFileSet()
		file, err1 := parser.ParseFile(fileset, "", strings.Join([]string{
			"package dummy",
			declaration,
		}, "\n"), 0)

		if err1 == nil {

			declaration_symbol := ""
			declaration_code   := ""

			for _, decl := range file.Decls {

				gen_decl, ok0 := decl.(*ast.GenDecl)

				if ok0 == true {

					if gen_decl.Tok == token.TYPE {

						for _, spec := range gen_decl.Specs {

							type_spec, ok1 := spec.(*ast.TypeSpec)

							if ok1 == true {

								struct_type, ok2 := type_spec.Type.(*ast.StructType)

								if ok2 == true {

									if type_spec.Name != nil {
										declaration_symbol = type_spec.Name.Name
									}

									buffer := bytes.Buffer{}
									printer.Fprint(&buffer, token.NewFileSet(), gen_decl)
									declaration_code = strings.TrimSpace(buffer.String())

									// explicitly referencing struct_type to keep symmetry with func parsing
									_ = struct_type

									break

								}

							}

						}

					}

				}

				if declaration_symbol != "" {
					break
				}

			}

			if declaration_symbol == symbol {

				_, ok3 := tool.contents[resolved]

				if ok3 == false {
					tool.contents[resolved] = make(map[string]requirement_specification)
				}

				tool.contents[resolved][symbol] = requirement_specification{
					File:        resolved,
					Type:        "struct",
					Declaration: declaration_code,
					Symbol:      declaration_symbol,
					Behavior:    behavior,
				}

				err2 := writeRequirements(tool)

				if err2 == nil {
					return fmt.Sprintf("requirements.DefineStruct: %s defined as %s", declaration_symbol, declaration_code), nil
				} else {
					return "", fmt.Errorf("requirements.DefineStruct: %s", err2.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineStruct: Invalid syntax, struct symbol \"%s\" must be the same as symbol \"%s\".", declaration_symbol, symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineStruct: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineStruct: %s", err0.Error())
	}

}

func (tool *Requirements) Search(path string, symbol string) (string, error) {
	// TODO
	return "", fmt.Errorf("requirements.Search: Not implemented yet")
}
