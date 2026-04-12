package tools

import "encoding/json"
import "fmt"
import "os"

func writeChangelog(tool *Changelog) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, "./exocomp-changelog.json")

		if err0 == nil {

			bytes, err1 := json.MarshalIndent(tool.contents, "", "\t")

			if err1 == nil {

				err2 := os.WriteFile(resolved, bytes, 0666)

				if err2 == nil {
					return nil
				} else {
					return fmt.Errorf("writeChangelog: %s", err2.Error())
				}

			} else {
				return fmt.Errorf("writeChangelog: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("writeChangelog: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("writeChangelog: Invalid tool playground")
	}

}
