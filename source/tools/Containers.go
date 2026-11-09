package tools

import "exocomp/schemas"
import utils_podman "exocomp/utils/podman"
import "crypto/sha256"
import "fmt"
import "path/filepath"
import "slices"
import "sort"
import "strings"

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

func (tool *Containers) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Create" {

			image, ok1 := arguments["image"].(string)
			name, ok2 := arguments["name"].(string)

			if ok1 == true && ok2 == true {

				command := ""
				network := "bridge"

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

func (tool *Containers) Create(image string, name string, command string, arguments []string, network string) (string, error) {

	if tool.isAllowedImage(image) == true {

		if tool.isValidName(name) == true {

			if network == "bridge" || network == "none" {

				if command == "" && len(arguments) > 0 {
					return "", fmt.Errorf("containers.Create: Invalid parameter \"arguments\" requires \"command\".")
				}

				if command == "" {
					command = "sleep"
					arguments = []string{"infinity"}
				} else if tool.isValidCommand(command) == false {
					return "", fmt.Errorf("containers.Create: Invalid command \"%s\".", command)
				}

				err1 := utils_podman.Create(image, name, command, arguments, network, tool.Sandbox, tool.label())

				if err1 == nil {
					return fmt.Sprintf("containers.Create: Container \"%s\" created from image \"%s\".", name, image), nil
				} else {
					return "", fmt.Errorf("containers.Create: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("containers.Create: Invalid network \"%s\".", network)
			}

		} else {
			return "", fmt.Errorf("containers.Create: Invalid container name \"%s\".", name)
		}

	} else {
		return "", fmt.Errorf("containers.Create: Invalid image \"%s\": Attempt to use unallowed image.", image)
	}

}

func (tool *Containers) Exec(name string, command string, arguments []string, workdir string) (string, error) {

	if tool.isValidName(name) == true {

		if tool.isValidCommand(command) == true {

			if strings.HasPrefix(workdir, "/") == true && tool.isValidCommand(workdir) == true {

				_, err1 := utils_podman.Inspect(name, tool.label())

				if err1 == nil {

					output, err2 := utils_podman.Exec(name, command, arguments, workdir)

					if err2 == nil {

						command_line := strings.Join(append([]string{command}, arguments...), " ")

						return strings.Join([]string{
							fmt.Sprintf("containers.Exec: podman exec %s %s", name, command_line),
							output,
						}, "\n"), nil

					} else {
						return "", fmt.Errorf("containers.Exec: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("containers.Exec: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("containers.Exec: Invalid workdir \"%s\".", workdir)
			}

		} else {
			return "", fmt.Errorf("containers.Exec: Invalid command \"%s\".", command)
		}

	} else {
		return "", fmt.Errorf("containers.Exec: Invalid container name \"%s\".", name)
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

func (tool *Containers) Images() (string, error) {

	images, err1 := utils_podman.Images()

	if err1 == nil {

		names := make([]string, 0)

		for _, image := range images {

			for _, tag := range image.RepoTags {

				if tool.isAllowedImage(tag) == true {
					names = append(names, tag)
				}

			}

		}

		sort.Strings(names)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("containers.Images: %d images available.", len(names)))

		for _, name := range names {
			result = append(result, fmt.Sprintf("Name: %s", name))
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("containers.Images: %s", err1.Error())
	}

}

func (tool *Containers) List() (string, error) {

	containers, err1 := utils_podman.Processes(tool.label())

	if err1 == nil {

		lines := make([]string, 0)

		for _, container := range containers {
			lines = append(lines, fmt.Sprintf("- Name: %s, Image: %s, Status: %s", strings.Join(container.Names, ","), container.Image, container.Status))
		}

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("containers.List: %d containers available.", len(lines)))

		for _, line := range lines {
			result = append(result, line)
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("containers.List: %s", err1.Error())
	}

}

func (tool *Containers) Logs(name string, tail int) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			output, err2 := utils_podman.Logs(name, tail)

			if err2 == nil {

				return strings.Join([]string{
					fmt.Sprintf("containers.Logs: %s", name),
					output,
				}, "\n"), nil

			} else {
				return "", fmt.Errorf("containers.Logs: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Logs: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Logs: Invalid container name \"%s\".", name)
	}

}

func (tool *Containers) Remove(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Remove(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Remove: Container \"%s\" removed.", name), nil
			} else {
				return "", fmt.Errorf("containers.Remove: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Remove: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Remove: Invalid container name \"%s\".", name)
	}

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

func (tool *Containers) Start(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Start(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Start: Container \"%s\" started.", name), nil
			} else {
				return "", fmt.Errorf("containers.Start: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Start: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Start: Invalid container name \"%s\".", name)
	}

}

func (tool *Containers) Stat(name string) (string, error) {

	if tool.isValidName(name) == true {

		container, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			id := container.ID

			if len(id) > 12 {
				id = id[0:12]
			}

			image := container.ImageName

			if image == "" {
				image = container.Config.Image
			}

			return strings.Join([]string{
				fmt.Sprintf("containers.Stat: %s", strings.TrimPrefix(container.Name, "/")),
				fmt.Sprintf("ID: %s", id),
				fmt.Sprintf("Image: %s", image),
				fmt.Sprintf("Status: %s", container.State.Status),
				fmt.Sprintf("Running: %t", container.State.Running),
				fmt.Sprintf("StartedAt: %s", container.State.StartedAt),
			}, "\n"), nil

		} else {
			return "", fmt.Errorf("containers.Stat: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Stat: Invalid container name \"%s\".", name)
	}

}

func (tool *Containers) Stop(name string) (string, error) {

	if tool.isValidName(name) == true {

		_, err1 := utils_podman.Inspect(name, tool.label())

		if err1 == nil {

			err2 := utils_podman.Stop(name)

			if err2 == nil {
				return fmt.Sprintf("containers.Stop: Container \"%s\" stopped.", name), nil
			} else {
				return "", fmt.Errorf("containers.Stop: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("containers.Stop: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("containers.Stop: Invalid container name \"%s\".", name)
	}

}

func (tool *Containers) isAllowedImage(image string) bool {

	if image == "" || len(image) > 255 {
		return false
	} else if strings.HasPrefix(image, "-") == true {
		return false
	}

	for _, character := range image {

		if character >= 'a' && character <= 'z' {
			continue
		} else if character >= 'A' && character <= 'Z' {
			continue
		} else if character >= '0' && character <= '9' {
			continue
		} else if strings.ContainsRune("._/:@-", character) == true {
			continue
		} else {
			return false
		}

	}

	return slices.Contains(tool.AllowedImages, image)

}

func (tool *Containers) isValidCommand(command string) bool {

	if command == "" {
		return false
	} else if strings.HasPrefix(command, "-") == true {
		return false
	} else if strings.ContainsAny(command, "\x00\n\r") == true {
		return false
	} else {
		return true
	}

}

func (tool *Containers) isValidName(name string) bool {

	if name == "" || len(name) > 64 {
		return false
	}

	for index, character := range name {

		if character >= 'a' && character <= 'z' {
			continue
		} else if character >= 'A' && character <= 'Z' {
			continue
		} else if character >= '0' && character <= '9' {
			continue
		} else if index > 0 && strings.ContainsRune("._-", character) == true {
			continue
		} else {
			return false
		}

	}

	return true

}

func (tool *Containers) label() string {

	sandbox, err0 := filepath.Abs(tool.Sandbox)

	if err0 == nil {

		resolved, err1 := filepath.EvalSymlinks(sandbox)

		if err1 == nil {
			sandbox = resolved
		}

	} else {
		sandbox = filepath.Clean(tool.Sandbox)
	}

	sum := sha256.Sum256([]byte(sandbox))

	return fmt.Sprintf("%x", sum[0:8])

}

func (tool *Containers) parseArguments(raw any) ([]string, error) {

	result := make([]string, 0)

	if raw == nil {
		return result, nil
	}

	values, ok := raw.([]interface{})

	if ok == true {

		result = make([]string, len(values))

		for index, value := range values {

			tmp, ok := value.(string)

			if ok == true {
				result[index] = tmp
			} else {
				return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")
			}

		}

		return result, nil

	} else {
		return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")
	}

}
