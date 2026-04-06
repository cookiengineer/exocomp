package tools

import "exocomp/utils"
import "fmt"

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

	} else if method == "DefineStruct" {

		// TODO

	} else if method == "DefineTest" {

		// TODO

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
}

func (tool *Requirements) DefineFunc(path string, symbol string, signature string, errors []string, description string) (string, error) {

}

func (tool *Requirements) DefineStruct(path string, symbol string, fields []string, description string) (string, error) {

}

func (tool *Requirements) DefineTest(path string, symbol string, inputs []string, outputs []string, description string) (string, error) {

}

func (tool *Requirements) Search(path string, symbol string) (string, error) {

}
