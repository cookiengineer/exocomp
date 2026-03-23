package gadgets

import "exocomp/utils"
import "fmt"
import "os"
import "sort"
import "strings"

type Files struct {
	Sandbox string
}

func NewFiles(sandbox string) *Files {

	return &Files{
		Sandbox: sandbox,
	}

}

func (gadget *Files) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!gadget:files.List \"./path/to/folder\"",
		"",
		"#!gadget:files.Read \"./path/to/file.go\"",
		"",
		"#!gadget:files.Stat \"./path/to/file.go\"",
		"",
		"#!gadget:files.Write \"./path/to/file.go\" <<#!EOF",
		"...file contents...",
		"#!EOF",
	}, "\n"), nil

}

func (gadget *Files) List(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(gadget.Sandbox, arguments[0])

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
						result = append(result, fmt.Sprintf("#!files.List: %s", resolved))

						for l := 0; l < len(lines); l++ {
							result = append(result, lines[l])
						}

						return strings.Join(result, "\n"), nil

					} else {
						return "", fmt.Errorf("#!files.List: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("#!files.List: \"%s\" is not a folder")
				}

			} else {
				return "", fmt.Errorf("#!files.List: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!files.List: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!files.List: Only one argument allowed")
	}

}

func (gadget *Files) Read(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(gadget.Sandbox, arguments[0])

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				result := strings.Join([]string{
					fmt.Sprintf("#!files.Read: %s", resolved),
					string(bytes),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("#!files.Read: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!files.Read: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!files.Read: Only one argument allowed")
	}

}

func (gadget *Files) Stat(arguments []string) (string, error) {

	if len(arguments) == 1 {

		resolved, err0 := resolveSandboxPath(gadget.Sandbox, arguments[0])

		if err0 == nil {

			stat, err1 := os.Stat(resolved)

			if err1 == nil {

				typ := "file"

				if stat.IsDir() == true {
					typ = "folder"
				}

				result := strings.Join([]string{
					fmt.Sprintf("#!files.Stat: %s", resolved),
					"Name: " + stat.Name(),
					"Type: " + typ,
					"Size: " + utils.FormatFileSize(stat.Size()),
					"Mode: " + utils.FormatFileMode(stat.Mode()),
					"Modified: " + utils.FormatTime(stat.ModTime()),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("#!files.Stat: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!files.Stat: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!files.Stat: Only one argument allowed")
	}

}

func (gadget *Files) Write(arguments []string) (string, error) {

	if len(arguments) == 2 {

		resolved, err0 := resolveSandboxPath(gadget.Sandbox, arguments[0])

		if err0 == nil {

			buffer, err1 := utils.FormatFileBuffer(arguments[1])

			if err1 == nil {

				err2 := os.WriteFile(resolved, buffer, 0666)

				if err2 == nil {

					result := strings.Join([]string{
						fmt.Sprintf("#!files.Write: %s", resolved),
						fmt.Sprintf("Buffer with %s written.", utils.FormatFileSize(int64(len(buffer)))),
					}, "\n")

					return result, nil

				} else {
					return "", fmt.Errorf("#!files.Write: %s", err2.Error())
				}

			} else {
				return "", fmt.Errorf("#!files.Write: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("#!files.Write: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("#!files.Write: Only two arguments allowed")
	}

}
