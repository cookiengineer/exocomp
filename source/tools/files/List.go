package files

import "fmt"
import "os"
import "sort"
import "strings"

func (tool *Files) List(path string) (string, error) {

	if path == "/" {
		path = "."
	} else if path == "" {
		path = "."
	}

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, false)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil {

			if stat.IsDir() == true {

				entries, err2 := os.ReadDir(resolved)

				if err2 == nil {

					lines := make([]string, 0)

					for _, entry := range entries {

						name := entry.Name()
						path := displayFilesPath(tool.Sandbox, resolved+"/"+name)

						if entry.Type()&os.ModeSymlink == 0 && strings.HasPrefix(name, ".") == false {

							typ := "file"

							if entry.IsDir() == true {
								typ = "folder"
							}

							lines = append(lines, fmt.Sprintf("- Path: \"%s\", Type: %s", path, typ))

						}

					}

					sort.Strings(lines)

					result := make([]string, 0)
					result = append(result, fmt.Sprintf("files.List: \"%s\" contains %d entries.", path, len(lines)))

					for l := 0; l < len(lines); l++ {
						result = append(result, lines[l])
					}

					return strings.Join(result, "\n"), nil

				} else {
					return "", fmt.Errorf("files.List: Cannot list folder \"%s\".", path)
				}

			} else {
				return "", fmt.Errorf("files.List: Invalid folder path \"%s\".", path)
			}

		} else {
			return sanitizeFilesystemError("files", "List", "Folder", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.List: %s", err0.Error())
	}

}
