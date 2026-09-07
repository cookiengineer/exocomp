package tools

import "exocomp/types"
import utils_ast "exocomp/utils/ast"
import "strings"

import "fmt"

func readImplementedRequirements(tool *Requirements) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, tool.Sandbox)

		if err0 == nil {

			package_symbols := utils_ast.GetPackageSymbols(resolved, false)

			for path, symbols := range package_symbols {

				path = strings.TrimPrefix(path, resolved + "/")

				tool.contents[path] = make(map[string]types.Requirement)

				for symbol, declaration := range symbols {

					tool.contents[path][symbol] = types.Requirement{
						Type:          declaration.Type,
						File:          path,
						Symbol:        symbol,
						Declaration:   declaration.Body,
						Behavior:      "",
						IsImplemented: true,
					}

				}

			}

			return nil

		} else {
			return fmt.Errorf("readImplementedRequirements: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readImplementedRequirements: Invalid tool playground")
	}

}
