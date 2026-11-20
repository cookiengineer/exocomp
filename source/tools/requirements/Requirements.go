package requirements

import "exocomp/schemas"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
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

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)
			behavior, ok4 := arguments["behavior"].(string)

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

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)
			behavior, ok4 := arguments["behavior"].(string)

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

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)
			behavior, ok4 := arguments["behavior"].(string)

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

		} else if method == "DefineType" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)
			behavior, ok4 := arguments["behavior"].(string)

			if ok1 == true && ok2 == true && ok3 == true && ok4 == true {
				return tool.DefineType(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(declaration), utils_fmt.FormatSingleLine(behavior))
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

			path, ok1 := arguments["path"].(string)
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

			path, ok1 := arguments["path"].(string)
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

	path := utils_fmt.FormatFilePath(id)
	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

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

func (tool *Requirements) MarshalJSON() ([]byte, error) {

	err := readRequirements(tool)

	if err == nil {
		return json.Marshal(tool.contents)
	} else {
		return nil, err
	}

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
