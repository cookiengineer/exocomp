package changelog

import "exocomp/schemas"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "sync"

type Changelog struct {
	Methods    []string
	Playground string
	Sandbox    string
	mutex      *sync.RWMutex
	contents   map[string]map[string][]types.ChangelogEntry // map[path][symbol]
}

func NewChangelog(methods []string, playground string, sandbox string) *Changelog {

	tool := &Changelog{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]map[string][]types.ChangelogEntry, 0),
	}

	return tool

}

func (tool *Changelog) Name() string {
	return "changelog"
}

func (tool *Changelog) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Add" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			description, ok3 := arguments["description"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Add(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Change" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			description, ok3 := arguments["description"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Change(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Deprecate" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			description, ok3 := arguments["description"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Deprecate(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Fix" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			description, ok3 := arguments["description"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Fix(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "List" {

			return tool.List()

		} else if method == "Remove" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			description, ok3 := arguments["description"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Remove(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Search" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)

			if ok1 == true && ok2 == true {
				return tool.Search(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
			} else if ok1 == true && ok2 == false {
				return tool.Search(utils_fmt.FormatFilePath(path), "")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
			}

		} else {
			return "", fmt.Errorf("changelog.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("changelog.%s: Method not allowed.", method)
	}

}

func (tool *Changelog) GetContent(id string) (any, error) {

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
				return nil, fmt.Errorf("changelog.GetContent: File \"%s\" has no Changelog entries.", path)
			}

		} else {
			return "", fmt.Errorf("changelog.GetContent: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.GetContent: %s", err1.Error())
	}

}

func (tool *Changelog) GetContentIdentifiers() []string {

	result := make([]string, 0)

	tool.mutex.RLock()

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	tool.mutex.RUnlock()

	sort.Strings(result)

	return result

}

func (tool *Changelog) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Changelog) MarshalJSON() ([]byte, error) {

	err := readChangelog(tool)

	if err == nil {
		return json.Marshal(tool.contents)
	} else {
		return nil, err
	}

}

func (tool *Changelog) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range ChangelogSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
