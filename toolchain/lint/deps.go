package main

import "strings"

func importPaths(file *File) []string {

	result := make([]string, 0, len(file.AST.Imports))

	for _, spec := range file.AST.Imports {

		value := importPath(spec)

		if value != "" {
			result = append(result, value)
		}

	}

	return result

}

func isToolPackage(relative string) (string, bool) {

	if strings.HasPrefix(relative, "tools/") == false {
		return "", false
	} else {

		parts := strings.Split(relative, "/")

		if len(parts) < 3 {
			return "", false
		} else {
			return parts[1], true
		}

	}

}

func checkToolIsolation(context *Context, file *File) []Diagnostic {

	current, ok := isToolPackage(file.Rel)

	if ok == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			prefix := context.Module + "/tools/"

			if strings.HasPrefix(imported, prefix) == true {

				other := strings.Split(strings.TrimPrefix(imported, prefix), "/")[0]

				if other != current {

					result = append(
						result,
						file.Diagnostic(
							"dep/tool-isolation",
							file.AST.Pos(),
							"tool package %s must not import tool package %s",
							current,
							other,
						),
					)

				}

			}

		}

		return result

	}

}

func checkToolsAllowed(context *Context, file *File) []Diagnostic {

	if _, ok := isToolPackage(file.Rel); ok == false || context.Module == "" {
		return nil
	} else {

		prefix := context.Module + "/"

		allowed := func(imported string) bool {

			if strings.HasPrefix(imported, prefix) == false {
				return true
			} else {

				rest := strings.TrimPrefix(imported, prefix)

				if rest == "types" || rest == "schemas" || rest == "agents" {
					return true
				} else {
					return strings.HasPrefix(rest, "utils/") == true || strings.HasPrefix(rest, "parsers/") == true
				}

			}

		}

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if allowed(imported) == false {

				result = append(
					result,
					file.Diagnostic(
						"dep/tools-allowed",
						file.AST.Pos(),
						"tool packages may only import types, schemas, agents and utils, got %q",
						imported,
					),
				)

			}

		}

		return result

	}

}

func checkTypesNoCapabilities(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "types/") == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if imported == context.Module+"/utils/capabilities" {

				result = append(
					result,
					file.Diagnostic(
						"dep/types-no-capabilities",
						file.AST.Pos(),
						"types must not import utils/capabilities",
					),
				)

			}

		}

		return result

	}

}

func checkUtilsNoTools(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "utils/") == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if strings.HasPrefix(imported, context.Module+"/tools") == true {

				result = append(
					result,
					file.Diagnostic(
						"dep/utils-no-tools",
						file.AST.Pos(),
						"utils must not import %q",
						imported,
					),
				)

			}

		}

		return result

	}

}
