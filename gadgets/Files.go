package gadgets

import "exocomp/types"
import "exocomp/utils"
import "fmt"
import "os"
import "strings"

type Files struct {
	Sandbox string
}

func NewFiles(config *types.Config) *Files {

	return &Files{
		Sandbox: config.Sandbox,
	}

}

func (gadget *Files) Help() string {

	return strings.Join([]string{
		"Files Gadget Usage:",
		"",
		"Read files with relative paths:",
		"",
		"#!gadget:files.Read \"./path/to/file.go\"",
		"",
		"Stat files with relative paths:",
		"",
		"#!gadget:files.Stat \"./path/to/file.go\"",
		"",
		"Write files with relative paths and heredoc syntax:",
		"",
		"#!gadget:files.Write \"./path/to/file.go\" <<#!EOF",
		"...file contents...",
		"#!EOF",
	}, "\n")

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

				result := strings.Join([]string{
					fmt.Sprintf("#!files.Stat: %s", resolved),
					"Name: " + stat.Name(),
					"Size: " + utils.FormatFileSize(stat.Size()),
					"Mode: " + utils.FormatFileMode(stat.Mode()),
					"Modified: " + utils.FormatTime(stat.ModTime()),
					"IsDirectory: " + utils.FormatBool(stat.IsDir()),
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
