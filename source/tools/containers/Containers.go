package containers

import "exocomp/schemas"
import "fmt"
import "slices"

type Containers struct {
	Methods       []string
	Playground    string
	Sandbox       string
	AllowedImages []string
}

func NewContainers(methods []string, playground string, sandbox string, allowed_images []string) *Containers {

	return &Containers{
		Methods:       methods,
		Playground:    playground,
		Sandbox:       sandbox,
		AllowedImages: allowed_images,
	}

}

func (tool *Containers) Name() string {
	return "containers"
}

func (tool *Containers) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Create" {

			image, ok1 := arguments["image"].(string)
			name, ok2 := arguments["name"].(string)

			if ok1 == true && ok2 == true {

				command := ""
				network := "none"

				raw_command, ok3 := arguments["command"]

				if ok3 == true {

					tmp, ok := raw_command.(string)

					if ok == true {
						command = tmp
					} else {
						return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"command\" is not a string.")
					}

				}

				raw_network, ok4 := arguments["network"]

				if ok4 == true {

					tmp, ok := raw_network.(string)

					if ok == true {
						network = tmp
					} else {
						return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"network\" is not a string.")
					}

				}

				args, err1 := tool.parseArguments(arguments["arguments"])

				if err1 == nil {
					return tool.Create(image, name, command, args, network)
				} else {
					return "", fmt.Errorf("containers.%s: %s", method, err1.Error())
				}

			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"image\" is not a string.")
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Exec" {

			name, ok1 := arguments["name"].(string)
			command, ok2 := arguments["command"].(string)

			if ok1 == true && ok2 == true {

				workdir := "/workspace"

				raw_workdir, ok3 := arguments["workdir"]

				if ok3 == true {

					tmp, ok := raw_workdir.(string)

					if ok == true {
						workdir = tmp
					} else {
						return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"workdir\" is not a string.")
					}

				}

				args, err1 := tool.parseArguments(arguments["arguments"])

				if err1 == nil {
					return tool.Exec(name, command, args, workdir)
				} else {
					return "", fmt.Errorf("containers.%s: %s", method, err1.Error())
				}

			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"command\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Images" {

			return tool.Images()

		} else if method == "List" {

			return tool.List()

		} else if method == "Logs" {

			name, ok1 := arguments["name"].(string)

			if ok1 == true {

				tail := 100

				raw_tail, ok2 := arguments["tail"]

				if ok2 == true {

					tmp, ok := raw_tail.(float64)

					if ok == true && tmp >= 1 && tmp <= 100000 {
						tail = int(tmp)
					} else {
						return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"tail\" is not a positive number.")
					}

				}

				return tool.Logs(name, tail)

			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Remove" {

			name, ok := arguments["name"].(string)

			if ok == true {
				return tool.Remove(name)
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Start" {

			name, ok := arguments["name"].(string)

			if ok == true {
				return tool.Start(name)
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Stat" {

			name, ok := arguments["name"].(string)

			if ok == true {
				return tool.Stat(name)
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Stop" {

			name, ok := arguments["name"].(string)

			if ok == true {
				return tool.Stop(name)
			} else {
				return "", fmt.Errorf("containers.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else {
			return "", fmt.Errorf("containers.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("containers.%s: Method not allowed.", method)
	}

}

func (tool *Containers) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Containers) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Containers) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Containers) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range ContainersSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
