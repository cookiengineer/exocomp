package files

import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "os"
import "strings"

func (tool *Files) Stat(path string) (string, error) {

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, false)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil {

			typ := "file"

			if stat.IsDir() == true {
				typ = "folder"
			}

			result := strings.Join([]string{
				fmt.Sprintf("files.Stat: \"%s\" is a %s.", path, typ),
				"Name: " + stat.Name(),
				"Type: " + typ,
				"Size: " + utils_fmt.FormatFileSize(stat.Size()),
				"Mode: " + utils_fmt.FormatFileMode(stat.Mode()),
				"Modified: " + utils_fmt.FormatTime(stat.ModTime()),
			}, "\n")

			return result, nil

		} else {
			return sanitizeFilesystemError("files", "Stat", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.Stat: %s", err0.Error())
	}

}
