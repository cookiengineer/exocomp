package tools

import "exocomp/types"
import utils_ast "exocomp/utils/ast"

import "fmt"

func readImplementedRequirements(tool *Requirements) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, tool.Sandbox)

		if err0 == nil {

			package_symbols := utils_ast.GetPackageSymbols(resolved, false)

			for path, symbols := range package_symbols {

				internal_path, err1 := sanitizeSandboxPath(tool.Playground, path)

				if err1 == nil {

					tool.contents[internal_path] = make(map[string]types.Requirement)

					for symbol, declaration := range symbols {

						tool.contents[internal_path][symbol] = types.Requirement{
							Type:          declaration.Type,
							File:          internal_path,
							Symbol:        symbol,
							Declaration:   declaration.Body,
							Behavior:      "",
							IsImplemented: true,
						}

					}

				} else {
					return fmt.Errorf("readImplementedRequirements: %s", err1.Error())
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
