package tools

import "exocomp/utils"
import "fmt"
import "os"
import "sort"
import "strings"

type Files struct {
	Sandbox string
}

func NewFiles(agent string, sandbox string) *Files {

	return &Files{
		Sandbox: sandbox,
	}

}

func (tool *Files) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		path, ok := arguments["path"].(string)

		if ok == true {
			return tool.List(path)
		} else {
			return "", fmt.Errorf("files.List: Invalid parameters")
		}

	} else if method == "Read" {

		path, ok := arguments["path"].(string)

		if ok == true {
			return tool.Read(path)
		} else {
			return "", fmt.Errorf("files.Read: Invalid parameters")
		}

	} else if method == "Stat" {

		path, ok := arguments["path"].(string)

		if ok == true {
			return tool.Stat(path)
		} else {
			return "", fmt.Errorf("files.Stat: Invalid parameters")
		}

	} else if method == "Write" {

		path,    ok1 := arguments["path"].(string)
		content, ok2 := arguments["content"].(string)

		if ok1 == true && ok2 == true {
			return tool.Write(path, content)
		} else {
			return "", fmt.Errorf("files.Write: Invalid parameters")
		}

	} else {
		return "", fmt.Errorf("files.%s: Invalid method.", method)
	}

}

func (tool *Files) List(path string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

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
					result = append(result, fmt.Sprintf("files.List: %s", resolved))

					for l := 0; l < len(lines); l++ {
						result = append(result, lines[l])
					}

					return strings.Join(result, "\n"), nil

				} else {
					return "", fmt.Errorf("files.List: %s", err2.Error())
				}

			} else {
				return "", fmt.Errorf("files.List: Invalid folder path \"%s\".")
			}

		} else {
			return "", fmt.Errorf("files.List: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("files.List: %s", err0.Error())
	}

}

func (tool *Files) Read(path string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		bytes, err1 := os.ReadFile(resolved)

		if err1 == nil {

			result := strings.Join([]string{
				fmt.Sprintf("files.Read: %s", resolved),
				string(bytes),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("files.Read: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("files.Read: %s", err0.Error())
	}

}

func (tool *Files) Stat(path string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil {

			typ := "file"

			if stat.IsDir() == true {
				typ = "folder"
			}

			result := strings.Join([]string{
				fmt.Sprintf("files.Stat: %s", resolved),
				"Name: " + stat.Name(),
				"Type: " + typ,
				"Size: " + utils.FormatFileSize(stat.Size()),
				"Mode: " + utils.FormatFileMode(stat.Mode()),
				"Modified: " + utils.FormatTime(stat.ModTime()),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("files.Stat: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("files.Stat: %s", err0.Error())
	}

}

func (tool *Files) Write(path string, content string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		buffer, err1 := utils.FormatFileBuffer(content)

		if err1 == nil {

			err2 := os.WriteFile(resolved, buffer, 0666)

			if err2 == nil {

				result := strings.Join([]string{
					fmt.Sprintf("files.Write: File %s with %s written.", resolved, utils.FormatFileSize(int64(len(buffer)))),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("files.Write: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("files.Write: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("files.Write: %s", err0.Error())
	}

}
