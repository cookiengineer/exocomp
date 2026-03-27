package tools

import "exocomp/utils"
import "fmt"
import "os"
import "sort"
import "strings"

type Files struct {
	Method    string
	Arguments []string
	Sandbox   string
}

func NewFiles(agent string, sandbox string, tools []string, programs []string) *Files {

	return &Files{
		Method:    "",
		Arguments: make([]string, 0),
		Sandbox:   sandbox,
	}

}

func (tool *Files) Call() (string, error) {

	if tool.Method == "Help" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "List" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "Read" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "Stat" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "Write" {
		return tool.List(tool.Arguments)
	} else {
		return "", fmt.Errorf("#!tool:files.%s: Invalid method.", tool.Method)
	}

}

func (tool *Files) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!tool:files.List \"./path/to/folder\"",
		"",
		"#!tool:files.Read \"./path/to/file.go\"",
		"",
		"#!tool:files.Stat \"./path/to/file.go\"",
		"",
		"#!tool:files.Write \"./path/to/file.go\" <<#!EOF",
		"...file contents...",
		"#!EOF",
	}, "\n"), nil

}

func (tool *Files) List(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, arguments[0])

		if err0 == nil {

			stat, err1 := os.Stat(resolved)

			if err1 == nil {

				if stat.IsDir() == true {

					entries, err2 := os.ReadDir(resolved)

					if err2 == nil {

						lines := make([]string, 0)

						for _, entry := range entries {

							name := entry.Name()

							if strings.HasPrefix(name, ".") == false {

								typ := "file"

								if entry.IsDir() == true {
									typ = "folder"
								}

								lines = append(lines, strings.Join([]string{
									"Name: " + name,
									"Type: " + typ,
								}, ", "))

							}

						}

						sort.Strings(lines)

						result := make([]string, 0)
						result = append(result, fmt.Sprintf("#!tool:files.List: %s", resolved))

						for l := 0; l < len(lines); l++ {
							result = append(result, lines[l])
						}

						return strings.Join(result, "\n"), nil

					} else {
						return "", fmt.Errorf("#!tool:files.List: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("#!tool:files.List: Invalid folder path \"%s\".")
				}

			} else {
				return "", fmt.Errorf("#!tool:files.List: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!tool:files.List: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!tool:files.List: Invalid arguments, only one argument allowed.")
	}

}

func (tool *Files) Parse(text string) (Tool, [2]int, error) {

	// #!tool:files.Read <path>
	// #!tool:files.Stat <path>
	// #!tool:files.Write <path> <<#!EOF
	// ...
	// #!EOF

	lines := strings.Split(text, "\n")

	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!tool:files.") {

		fields := utils.SplitArguments(strings.TrimSpace(lines[0][len("#!tool:files."):]))
		method := strings.ToUpper(fields[0][0:1]) + strings.ToLower(fields[0][1:])
		parsed := [2]int{0, 1}

		for f := 1; f < len(fields); f++ {

			field := fields[f]

			if strings.HasPrefix(field, "<<") {

				heredoc_marker := field[2:]

				for s := 1; s < len(lines); s++ {

					if strings.HasPrefix(lines[s], heredoc_marker) {

						fields[f] = strings.Join(lines[1:s], "\n")
						parsed[1] = int(s)

						break

					}

				}

			}

		}

		if method == "Help" ||
			method == "List" ||
			method == "Read" ||
			method == "Stat" ||
			method == "Write" {

			tool.Method    = method
			tool.Arguments = fields[1:]

			return Tool(tool), parsed, nil

		} else {
			return nil, [2]int{0, len(lines)}, fmt.Errorf("#!tool:files.%s: Invalid method.", method)
		}

	} else {
		return nil, [2]int{0, len(lines)}, fmt.Errorf("Invalid Tool Call line.")
	}

}

func (tool *Files) Read(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, arguments[0])

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				result := strings.Join([]string{
					fmt.Sprintf("#!tool:files.Read: %s", resolved),
					string(bytes),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("#!tool:files.Read: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!tool:files.Read: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!tool:files.Read: Invalid arguments, only one argument allowed.")
	}

}

func (tool *Files) Stat(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, arguments[0])

		if err0 == nil {

			stat, err1 := os.Stat(resolved)

			if err1 == nil {

				typ := "file"

				if stat.IsDir() == true {
					typ = "folder"
				}

				result := strings.Join([]string{
					fmt.Sprintf("#!tool:files.Stat: %s", resolved),
					"Name: " + stat.Name(),
					"Type: " + typ,
					"Size: " + utils.FormatFileSize(stat.Size()),
					"Mode: " + utils.FormatFileMode(stat.Mode()),
					"Modified: " + utils.FormatTime(stat.ModTime()),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("#!tool:files.Stat: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!tool:files.Stat: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!tool:files.Stat: Invalid arguments, only one argument allowed.")
	}

}

func (tool *Files) Write(arguments []string) (string, error) {

	if len(arguments) == 2 {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, arguments[0])

		if err0 == nil {

			buffer, err1 := utils.FormatFileBuffer(arguments[1])

			if err1 == nil {

				err2 := os.WriteFile(resolved, buffer, 0666)

				if err2 == nil {

					result := strings.Join([]string{
						fmt.Sprintf("#!tool:files.Write: File %s with %s written.", resolved, utils.FormatFileSize(int64(len(buffer)))),
					}, "\n")

					return result, nil

				} else {
					return "", fmt.Errorf("#!tool:files.Write: %s", err2.Error())
				}

			} else {
				return "", fmt.Errorf("#!tool:files.Write: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!tool:files.Write: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!tool:files.Write: Invalid arguments, only two arguments allowed.")
	}

}
