package tools

import "exocomp/schemas"
import utils_ast "exocomp/utils/ast"
import utils_fmt "exocomp/utils/fmt"
import utils_fs "exocomp/utils/fs"
import "bytes"
import "fmt"
import "os"
import "path/filepath"
import "slices"
import "sort"
import "strings"

type Files struct {
	Methods []string
	Sandbox string
}

func NewFiles(methods []string, playground string, sandbox string) *Files {

	return &Files{
		Methods: methods,
		Sandbox: sandbox,
	}

}

func (tool *Files) Name() string {
	return "files"
}

func (tool *Files) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Copy" {

			from_path, ok1 := arguments["from_path"].(string)
			to_path,   ok2 := arguments["to_path"].(string)

			if ok1 == true && ok2 == true {
				return tool.Copy(utils_fmt.FormatFilePath(from_path), utils_fmt.FormatFilePath(to_path))
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"to_path\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"from_path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "List" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.List(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Read" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.Read(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Stat" {

			path, ok := arguments["path"].(string)

			if ok == true {
				return tool.Stat(utils_fmt.FormatFilePath(path))
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			}

		} else if method == "Search" {

			path,  ok1 := arguments["path"].(string)
			query, ok2 := arguments["query"].(string)

			if ok1 == true && ok2 == true {
				return tool.Search(utils_fmt.FormatFilePath(path), query)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"query\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "ReadSymbol" {

			path,   ok1 := arguments["path"].(string)
			symbol, ok2 := arguments["symbol"].(string)

			if ok1 == true && ok2 == true {
				return tool.ReadSymbol(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "WriteSymbol" {

			path,        ok1 := arguments["path"].(string)
			symbol,      ok2 := arguments["symbol"].(string)
			declaration, ok3 := arguments["declaration"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.WriteSymbol(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatMultiLine(declaration))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"declaration\" is not a string.")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Write" {

			path,    ok1 := arguments["path"].(string)
			content, ok2 := arguments["content"].(string)

			if ok1 == true && ok2 == true {
				return tool.Write(utils_fmt.FormatFilePath(path), content)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"content\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameter \"path\" is not a string.")
			} else {
				return "", fmt.Errorf("files.%s: %s", method, "Invalid parameters.")
			}

		} else {
			return "", fmt.Errorf("files.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("files.%s: Method not allowed.", method)
	}

}

func (tool *Files) Copy(from_path string, to_path string) (string, error) {

	from_resolved, err1 := resolveSandboxPath(tool.Sandbox, from_path)

	if err1 == nil {

		to_resolved, err2 := resolveSandboxPath(tool.Sandbox, to_path)

		if err2 == nil {

			stat, err3 := os.Stat(from_resolved)

			if err3 == nil {

				if stat.IsDir() {

					err4 := utils_fs.CopyAll(from_resolved, to_resolved)

					if err4 == nil {
						return fmt.Sprintf("files.Copy: Folder \"%s\" copied to \"%s\".", from_path, to_path), nil
					} else {
						return "", fmt.Errorf("files.Copy: Cannot copy \"%s\" to \"%s\".", from_path, to_path)
					}

				} else {

					err4 := utils_fs.Copy(from_resolved, to_resolved)

					if err4 == nil {
						return fmt.Sprintf("files.Copy: File \"%s\" copied to \"%s\".", from_path, to_path), nil
					} else {
						return "", fmt.Errorf("files.Copy: Cannot copy \"%s\" to \"%s\".", from_path, to_path)
					}

				}

			} else {
				return sanitizeFilesystemError("files", "Copy", "Path", from_path, err3)
			}

		} else {
			return "", fmt.Errorf("files.Copy: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("files.Copy: %s", err1.Error())
	}

}

func (tool *Files) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Files) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Files) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Files) List(path string) (string, error) {

	if path == "/" {
		path = "."
	} else if path == "" {
		path = "."
	}

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil {

			if stat.IsDir() == true {

				entries, err2 := os.ReadDir(resolved)

				if err2 == nil {

					lines := make([]string, 0)

					for _, entry := range entries {

						name      := entry.Name()
						path, err := sanitizeSandboxPath(tool.Sandbox, resolved + "/" + name)

						if err == nil && strings.HasPrefix(name, ".") == false {

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

func (tool *Files) Read(path string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

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

func (tool *Files) Search(path string, query string) (string, error) {

	if path == "/" {
		path = "."
	} else if path == "" {
		path = "."
	}

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil {

			if stat.IsDir() == true {

				lines := make([]string, 0)

				err2 := filepath.WalkDir(resolved, func(entry_path string, entry os.DirEntry, err3 error) error {

					if err3 == nil {

						name := entry.Name()

						if strings.HasPrefix(name, ".") == false {

							if entry.IsDir() == false && strings.HasSuffix(name, ".go") {

								bytes, err4 := os.ReadFile(entry_path)

								if err4 == nil {

									symbols := utils_ast.SearchSymbols(bytes, query)

									for s := 0; s < len(symbols); s++ {

										sandbox_path, err5 := sanitizeSandboxPath(tool.Sandbox, entry_path)

										if err5 == nil {
											lines = append(lines, fmt.Sprintf("- File: \"%s\", Symbol: \"%s\", Type: \"%s\"", sandbox_path, symbols[s].Name, symbols[s].Type))
										}

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

func (tool *Files) ReadSymbol(path string, symbol string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		bytes, err1 := os.ReadFile(resolved)

		if err1 == nil {

			found := utils_ast.GetSymbol(bytes, symbol, "")

			if found != nil {

				result := strings.Join([]string{
					fmt.Sprintf("files.ReadSymbol: File \"%s\" Symbol \"%s\" with Type \"%s\"", path, found.Name, found.Type),
					fmt.Sprintf("%s", found.Body),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("files.ReadSymbol: File \"%s\" has no Symbol \"%s\"", path, symbol)
			}

		} else {
			return sanitizeFilesystemError("files", "ReadSymbol", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.ReadSymbol: %s", err0.Error())
	}

}

func (tool *Files) WriteSymbol(path string, symbol string, declaration string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		source, err1 := os.ReadFile(resolved)

		if err1 == nil {

			declaration_type := ""

			if strings.HasPrefix(strings.TrimSpace(declaration), "func") {
				declaration_type = "func"
			}

			result := utils_ast.WriteSymbol(source, symbol, declaration, declaration_type)

			if bytes.Equal(result, source) == true {
				return "", fmt.Errorf("files.WriteSymbol: File \"%s\" has no Symbol \"%s\"", path, symbol)
			}

			err2 := os.WriteFile(resolved, result, 0666)

			if err2 == nil {
				return fmt.Sprintf("files.WriteSymbol: File \"%s\" Symbol \"%s\" written.", path, symbol), nil
			} else {
				return sanitizeFilesystemError("files", "WriteSymbol", "File", path, err2)
			}

		} else {
			return sanitizeFilesystemError("files", "WriteSymbol", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.WriteSymbol: %s", err0.Error())
	}

}

func (tool *Files) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range FilesSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

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

func (tool *Files) Write(path string, content string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, path)

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
