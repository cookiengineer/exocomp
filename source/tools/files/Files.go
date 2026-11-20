package files

import "exocomp/schemas"
import utils_capabilities "exocomp/utils/capabilities"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "slices"

type Files struct {
	Methods      []string
	Sandbox      string
	Capabilities *utils_capabilities.Set
}

func NewFiles(methods []string, playground string, sandbox string) *Files {

	return NewFilesWithCapabilities(methods, playground, sandbox, nil)

}

func NewFilesWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	capabilities *utils_capabilities.Set,
) *Files {

	return &Files{
		Methods:      methods,
		Sandbox:      sandbox,
		Capabilities: capabilities,
	}

}

func (tool *Files) Name() string {
	return "files"
}

func (tool *Files) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Copy" {

			from_path, ok1 := arguments["from_path"].(string)
			to_path, ok2 := arguments["to_path"].(string)

			if ok1 == true && ok2 == true {
				return tool.Copy(utils_fmt.FormatFilePath(from_path), utils_fmt.FormatFilePath(to_path))
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"to_path\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"from_path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "List" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.List(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Read" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.Read(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Stat" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.Stat(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Search" {

			path, ok1 := arguments["path"].(string)
			query, ok2 := arguments["query"].(string)

			if ok1 == true && ok2 == true {
				return tool.Search(utils_fmt.FormatFilePath(path), query)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"query\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "ReadSymbol" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)

			if ok1 == true && ok2 == true {
				return tool.ReadSymbol(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "WriteSymbol" {

			path, ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.WriteSymbol(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatMultiLine(declaration))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"declaration\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Write" {

			path, ok1 := arguments["path"].(string)
			content, ok2 := arguments["content"].(string)

			if ok1 == true && ok2 == true {
				return tool.Write(utils_fmt.FormatFilePath(path), content)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"content\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else {
			return "", fmt.Errorf("files.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("files.%s: Method not allowed.", method)
	}

}

func (tool *Files) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Files) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Files) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Files) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range FilesSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
