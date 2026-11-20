package files

import utils_ast "exocomp/utils/ast"
import "fmt"
import "os"
import "path/filepath"
import "sort"
import "strings"

func (tool *Files) Search(path string, query string) (string, error) {

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

				lines := make([]string, 0)

				err2 := filepath.WalkDir(resolved, func(entry_path string, entry os.DirEntry, err3 error) error {

					if err3 == nil {

						name := entry.Name()

						if entry.Type()&os.ModeSymlink != 0 {
							return nil
						}

						if strings.HasPrefix(name, ".") == false {

							if entry.IsDir() == false && strings.HasSuffix(name, ".go") {

								bytes, err4 := os.ReadFile(entry_path)

								if err4 == nil {

									symbols := utils_ast.SearchSymbols(bytes, query)

									for s := 0; s < len(symbols); s++ {

										sandbox_path := displayFilesPath(tool.Sandbox, entry_path)
										lines = append(lines, fmt.Sprintf("- File: \"%s\", Symbol: \"%s\", Type: \"%s\"", sandbox_path, symbols[s].Name, symbols[s].Type))

									}

								}

							}

						} else if entry.IsDir() == true {
							return filepath.SkipDir
						}

					}

					return nil

				})

				if err2 == nil {

					sort.Strings(lines)

					result := make([]string, 0)
					result = append(result, fmt.Sprintf("files.Search: \"%s\" for \"%s\" contains %d matches.", path, query, len(lines)))

					for l := 0; l < len(lines); l++ {
						result = append(result, lines[l])
					}

					return strings.Join(result, "\n"), nil

				} else {
					return "", fmt.Errorf("files.Search: Cannot search folder \"%s\".", path)
				}

			} else {
				return "", fmt.Errorf("files.Search: Invalid folder path \"%s\".", path)
			}

		} else {
			return sanitizeFilesystemError("files", "Search", "Folder", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.Search: %s", err0.Error())
	}

}
