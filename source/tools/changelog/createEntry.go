package changelog

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "time"

func (tool *Changelog) createEntry(method string, path string, symbol string, description string) (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readChangelog(tool)

			_, ok1 := tool.contents[internal_path]

			if ok1 == false {
				tool.contents[internal_path] = make(map[string][]types.ChangelogEntry, 0)
			}

			_, ok2 := tool.contents[internal_path][symbol]

			if ok2 == false {
				tool.contents[internal_path][symbol] = make([]types.ChangelogEntry, 0)
			}

			var found *types.ChangelogEntry = nil

			for _, entry := range tool.contents[internal_path][symbol] {

				if entry.Type == method && entry.File == internal_path && entry.Symbol == symbol && entry.Description == description {

					found = &entry
					break

				}

			}

			if found == nil {

				now := time.Now()
				today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

				tool.contents[internal_path][symbol] = append(tool.contents[internal_path][symbol], types.ChangelogEntry{
					Date:        today,
					Type:        method,
					File:        internal_path,
					Symbol:      symbol,
					Description: description,
				})

				err3 := writeChangelog(tool)

				if err3 == nil {
					return fmt.Sprintf("changelog.%s: Log entry created for File \"%s\" with Symbol \"%s\" at %s.", method, path, symbol, today.Format("2006-01-02")), nil
				} else {
					return "", fmt.Errorf("changelog.%s: %s", method, err3.Error())
				}

			} else {
				return fmt.Sprintf("changelog.%s: Log entry already exists for File \"%s\" with Symbol \"%s\" at %s.", method, path, found.Symbol, found.Date.Format("2006-01-02")), nil
			}

		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.%s: %s", method, err1.Error())
	}

}
