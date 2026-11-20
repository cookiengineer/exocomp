package files

import "fmt"
import "os"
import "strings"

func (tool *Files) Read(path string) (string, error) {

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, false)

	if err0 == nil {

		bytes, err1 := os.ReadFile(resolved)

		if err1 == nil {

			result := strings.Join([]string{
				fmt.Sprintf("files.Read: File \"%s\" contents", path),
				string(bytes),
			}, "\n")

			return result, nil

		} else {
			return sanitizeFilesystemError("files", "Read", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.Read: %s", err0.Error())
	}

}
