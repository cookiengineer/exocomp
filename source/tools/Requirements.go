package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import utils_ast "exocomp/utils/ast"
import "encoding/json"
import "fmt"
import "os"
import "slices"
import "sort"
import "strings"
import "sync"

type Requirements struct {
	Methods    []string
	Playground string
	Sandbox    string
	mutex      *sync.RWMutex
	contents   map[string]map[string]types.Requirement // map[resolved][symbol]
}

func NewRequirements(methods []string, playground string, sandbox string) *Requirements {

	tool := &Requirements{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]map[string]types.Requirement),
	}

	readRequirements(tool)

	return tool

}

func (tool *Requirements) Name() string {
	return "requirements"
}

func (tool *Requirements) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

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

		} else if method == "DefineInterface" {

			path,        ok1 := arguments["path"].(string)
			symbol,      ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)
			behavior,    ok4 := arguments["behavior"].(string)

			if ok1 == true && ok2 == true && ok3 == true && ok4 == true {
				return tool.DefineInterface(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatMultiLine(declaration), utils_fmt.FormatSingleLine(behavior))
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
				return tool.DefineStruct(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatMultiLine(declaration), utils_fmt.FormatSingleLine(behavior))
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

		} else if method == "Signoff" {

			path,   ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)

			if ok1 == true && ok2 == true {
				return tool.Signoff(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("requirements.%s: %s", method, "Invalid parameters.")
			}

		} else {
			return "", fmt.Errorf("requirements.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("requirements.%s: Method not allowed.", method)
	}

}

func (tool *Requirements) GetContent(id string) (any, error) {

	path       := utils_fmt.FormatFilePath(id)
	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			tool.mutex.RLock()
			content, ok := tool.contents[internal_path]
			tool.mutex.RUnlock()

			if ok == true {
				return content, nil
			} else {
				return nil, fmt.Errorf("requirements.GetContent: No specification defined for path \"%s\".", path)
			}

		} else {
			return "", fmt.Errorf("requirements.GetContent: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.GetContent: %s", err1.Error())
	}

}

func (tool *Requirements) GetContentIdentifiers() []string {

	result := make([]string, 0)

	tool.mutex.RLock()

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	tool.mutex.RUnlock()

	sort.Strings(result)

	return result

}

func (tool *Requirements) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Requirements) List() (string, error) {

	readRequirements(tool)

	lines := make([]string, 0)

	tool.mutex.RLock()

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

	tool.mutex.RUnlock()

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("requirements.List: %d specifications.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Requirements) MarshalJSON() ([]byte, error) {

	err := readRequirements(tool)

	if err != nil {
		return nil, err
	}

	return json.Marshal(tool.contents)

}

func (tool *Requirements) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range RequirementsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}

func (tool *Requirements) Signoff(path string, symbol string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			readRequirements(tool)

			tool.mutex.RLock()
			specification, ok1 := tool.contents[internal_path][symbol]
			tool.mutex.RUnlock()

			if ok1 == true {

				resolved_path, err3 := resolveSandboxPath(tool.Playground, specification.File)

				if err3 == nil {

					source, err4 := os.ReadFile(resolved_path)

					if err4 == nil {

						if utils_ast.HasSymbol(source, specification.Symbol, specification.Type) == true {

							tool.mutex.Lock()
							specification.IsImplemented = true
							tool.contents[internal_path][symbol] = specification
							tool.mutex.Unlock()

							err5 := writeRequirements(tool)

							if err5 == nil {
								return fmt.Sprintf("requirements.Signoff: %s#%s marked as implemented.", path, symbol), nil
							} else {
								return "", fmt.Errorf("requirements.Signoff: %s", err5.Error())
							}

						} else {
							return "", fmt.Errorf("requirements.Signoff: Symbol \"%s\" is not implemented in \"%s\" yet.", symbol, path)
						}

					} else {
						return "", fmt.Errorf("requirements.Signoff: %s", err4.Error())
					}

				} else {
					return "", fmt.Errorf("requirements.Signoff: %s", err3.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.Signoff: No specification available for path \"%s\" and symbol \"%s\"", path, symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.Signoff: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.Signoff: %s", err1.Error())
	}

}

func (tool *Requirements) DefineFunc(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "func")

			if prettified != "" {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok1 := tool.contents[internal_path]

				if ok1 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "func",
					Declaration: prettified,
					Symbol:      symbol,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineFunc: %s defined as %s", symbol, prettified), nil
				} else {
					return "", fmt.Errorf("requirements.DefineFunc: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineFunc: Invalid Go syntax. \"func %s (...) (...)\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineFunc: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineFunc: %s", err1.Error())
	}

}

func (tool *Requirements) DefineInterface(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			if strings.HasPrefix(declaration, "interface{") || strings.HasPrefix(declaration, "interface {") {
				declaration = "type " + symbol + " interface " + strings.TrimSpace(declaration[9:])
			}

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "interface")

			if prettified != "" {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok3 := tool.contents[internal_path]

				if ok3 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "interface",
					Declaration: prettified,
					Symbol:      symbol,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineInterface: %s defined as %s", symbol, prettified), nil
				} else {
					return "", fmt.Errorf("requirements.DefineInterface: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineInterface: Invalid Go syntax. \"type %s interface { ...}\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineInterface: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineInterface: %s", err1.Error())
	}

}

func (tool *Requirements) DefineStruct(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			if strings.HasPrefix(declaration, "struct{") || strings.HasPrefix(declaration, "struct {") {
				declaration = "type " + symbol + " struct " + strings.TrimSpace(declaration[6:])
			}

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "struct")

			if prettified != "" {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok3 := tool.contents[internal_path]

				if ok3 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "struct",
					Declaration: prettified,
					Symbol:      symbol,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineStruct: %s defined as %s", symbol, prettified), nil
				} else {
					return "", fmt.Errorf("requirements.DefineStruct: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineStruct: Invalid Go syntax. \"type %s struct { ...}\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineStruct: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineStruct: %s", err1.Error())
	}

}

func (tool *Requirements) Search(path string, symbol string) (string, error) {
	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		search_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			readRequirements(tool)

			lines := make([]string, 0)

			tool.mutex.RLock()

			for internal_path, specifications := range tool.contents {

				if strings.HasPrefix(internal_path, search_path) {

					for internal_symbol, specification := range specifications {

						if symbol != "" {

							matches_symbol := strings.Contains(strings.ToLower(internal_symbol), strings.ToLower(symbol))

							if matches_symbol == true {

								sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, specification.File)

								if err3 == nil {
									lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Declaration: %s, Behavior: %s", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior))
								}

							}

						} else if symbol == "" {

							sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, specification.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Declaration: %s, Behavior: %s", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior))
							}

						}

					}

				}

			}

			tool.mutex.RUnlock()

			sort.Strings(lines)

			result := make([]string, 0)

			if symbol != "" {
				result = append(result, fmt.Sprintf("requirements.Search: %s#%s* contains %d specifications.", path, symbol, len(lines)))
			} else {
				result = append(result, fmt.Sprintf("requirements.Search: %s#* contains %d specifications.", path, len(lines)))
			}

			for l := 0; l < len(lines); l++ {
				result = append(result, lines[l])
			}

			return strings.Join(result, "\n"), nil

		} else {
			return "", fmt.Errorf("requirements.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.Search: %s", err1.Error())
	}

}

