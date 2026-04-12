package tools

import "encoding/json"
import "fmt"
import "os"

func readChangelog(tool *Changelog) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, "./exocomp-changelog.json")

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				contents := make(map[string]map[string][]changelog_entry)
				err2     := json.Unmarshal(bytes, &contents)

				if err2 == nil {

					for file, _ := range tool.contents {
						delete(tool.contents, file)
					}

					for file, symbols := range contents {
						tool.contents[file] = symbols
					}

					return nil

				} else {
					return fmt.Errorf("readChangelog: %s", err2.Error())
				}

			} else {
				return fmt.Errorf("readChangelog: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readChangelog: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readChangelog: Invalid tool playground")
	}

}
