package tools

import "exocomp/types"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func readBugs(tool *Bugs) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, filepath.Join(".exocomp", "bugs.json"))

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				contents := make(map[string]map[string]types.Bug)
				err2 := json.Unmarshal(bytes, &contents)

				if err2 == nil {

					for file, _ := range tool.contents {
						delete(tool.contents, file)
					}

					for file, symbols := range contents {
						tool.contents[file] = symbols
					}

					return nil

				} else {
					return fmt.Errorf("bugs.Read: Cannot parse bug reports.")
				}

			} else if os.IsNotExist(err1) {

				tool.contents = make(map[string]map[string]types.Bug)
				return nil

			} else {
				return fmt.Errorf("bugs.Read: Cannot read bug reports.")
			}

		} else {
			return fmt.Errorf("bugs.Read: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("bugs.Read: Invalid tool playground.")
	}

}
