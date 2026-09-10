package tools

import "exocomp/types"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func writeRequirements(tool *Requirements) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, filepath.Join(".exocomp", "requirements.json"))

		if err0 == nil {

			contents := make(map[string]map[string]types.Requirement)

			for resolved, specifications := range tool.contents {

				for symbol, specification := range specifications {

					if specification.IsImplemented == false {

						_, ok := contents[resolved]

						if ok == false {
							contents[resolved] = make(map[string]types.Requirement)
						}

						contents[resolved][symbol] = specification
					}

				}

			}

			bytes, err1 := json.MarshalIndent(contents, "", "\t")

			if err1 == nil {

				err2 := os.WriteFile(resolved, bytes, 0666)

				if err2 == nil {
					return nil
				} else {
					return fmt.Errorf("writeRequirements: Cannot persist requirements")
				}

			} else {
				return fmt.Errorf("writeRequirements: Cannot serialize requirements")
			}

		} else {
			return fmt.Errorf("writeRequirements: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("writeRequirements: Invalid tool playground")
	}

}
