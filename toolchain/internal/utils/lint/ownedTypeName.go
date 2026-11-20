package lint

import "go/ast"
import "strings"

// ownedTypeName returns the type a free function belongs to, if any.
// Constructors, parsers and converters live in their type's file.
func ownedTypeName(func_decl *ast.FuncDecl, types map[string]*File) string {

	name := func_decl.Name.Name

	for _, prefix := range []string{"New", "Parse", "new"} {

		if strings.HasPrefix(name, prefix) == false {
			continue
		}

		rest := strings.TrimPrefix(name, prefix)
		longest := ""

		for type_name := range types {

			if rest == type_name || strings.EqualFold(rest, type_name) == true || strings.HasPrefix(rest, type_name) == true {

				if len(type_name) > len(longest) {
					longest = type_name
				}

			}

		}

		if longest != "" {
			return longest
		}

	}

	longest := ""

	for type_name := range types {

		if strings.HasPrefix(name, type_name+"From") == true {

			if len(type_name) > len(longest) {
				longest = type_name
			}

		}

	}

	if longest != "" {
		return longest
	} else if strings.HasPrefix(name, "New") == true || strings.HasPrefix(name, "Parse") == true || strings.HasPrefix(name, "new") == true {
		return declaredResultType(func_decl, types)
	} else {
		return ""
	}

}
