package programs

import "exocomp/schemas"
import utils_capabilities "exocomp/utils/capabilities"
import "fmt"
import "slices"

type Programs struct {
	Methods         []string
	Sandbox         string
	AllowedPrograms []string
	Capabilities    *utils_capabilities.Set
}

func NewPrograms(methods []string, playground string, sandbox string, allowed_programs []string) *Programs {

	return NewProgramsWithCapabilities(methods, playground, sandbox, allowed_programs, nil)

}

func NewProgramsWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	allowed_programs []string,
	capabilities *utils_capabilities.Set,
) *Programs {

	return &Programs{
		Methods:         methods,
		Sandbox:         sandbox,
		AllowedPrograms: allowed_programs,
		Capabilities:    capabilities,
	}

}

func (tool *Programs) Name() string {
	return "programs"
}

func (tool *Programs) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "List" {

			return tool.List()

		} else if method == "Execute" {

			program, ok1 := arguments["program"].(string)

			if ok1 == true {

				args := make([]string, 0)

				raw, ok2 := arguments["arguments"]

				if ok2 == true {

					raw_args, ok3 := raw.([]any)

					if ok3 == true {

						args = make([]string, len(raw_args))

						for a, value := range raw_args {

							tmp, ok := value.(string)

							if ok == true {
								args[a] = tmp
							}

						}

					} else {
						return "", fmt.Errorf("programs.%s: %s", method, "Invalid parameter \"arguments\" is not an array of strings.")
					}

				}

				return tool.Execute(program, args)

			} else {
				return "", fmt.Errorf("programs.%s: %s", method, "Invalid parameter \"program\" is not a string.")
			}

		} else if method == "Stat" {

			program, ok1 := arguments["program"].(string)

			if ok1 == true {
				return tool.Stat(program)
			} else {
				return "", fmt.Errorf("programs.%s: %s", method, "Invalid parameter \"program\" is not a string.")
			}

		} else {
			return "", fmt.Errorf("programs.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("programs.%s: Method not allowed.", method)
	}

}

func (tool *Programs) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Programs) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Programs) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Programs) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range ProgramsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
