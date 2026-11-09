package tools

import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func writeChangelog(tool *Changelog) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, filepath.Join(".exocomp", "changelog.json"))

		if err0 == nil {

			bytes, err1 := json.MarshalIndent(tool.contents, "", "\t")

			if err1 == nil {

				err2 := os.WriteFile(resolved, bytes, 0666)

				if err2 == nil {
					return nil
				} else {
					return fmt.Errorf("changelog.Write: Cannot persist changelog.")
				}

			} else {
				return fmt.Errorf("changelog.Write: Cannot serialize changelog.")
			}

		} else {
			return fmt.Errorf("changelog.Write: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("changelog.Write: Invalid tool playground.")
	}

}
