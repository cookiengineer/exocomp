package files

import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "os"
import "strings"

func (tool *Files) Write(path string, content string) (string, error) {

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, true)

	if err0 == nil {

		buffer, err1 := utils_fmt.FormatFileBuffer(content)

		if err1 == nil {

			err2 := os.WriteFile(resolved, buffer, 0666)

			if err2 == nil {

				result := strings.Join([]string{
					fmt.Sprintf("files.Write: File \"%s\" with %s written.", path, utils_fmt.FormatFileSize(int64(len(buffer)))),
				}, "\n")

				return result, nil

			} else {
				return sanitizeFilesystemError("files", "Write", "Folder", path, err2)
			}

		} else {
			return "", fmt.Errorf("files.Write: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("files.Write: %s", err0.Error())
	}

}
