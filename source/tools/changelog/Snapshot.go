package changelog

import "exocomp/types"
import "fmt"
import "sort"

func (tool *Changelog) Snapshot() ([]types.ChangelogEntry, error) {

	tool.mutex.Lock()
	err := readChangelog(tool)

	if err == nil {

		result := make([]types.ChangelogEntry, 0)

		for _, symbols := range tool.contents {

			for _, entries := range symbols {
				result = append(result, entries...)
			}

		}

		tool.mutex.Unlock()

		sort.Slice(result, func(a int, b int) bool {

			if result[a].Date.Equal(result[b].Date) == false {
				return result[a].Date.Before(result[b].Date)
			} else if result[a].File != result[b].File {
				return result[a].File < result[b].File
			} else if result[a].Symbol != result[b].Symbol {
				return result[a].Symbol < result[b].Symbol
			} else if result[a].Type != result[b].Type {
				return result[a].Type < result[b].Type
			} else {
				return result[a].Description < result[b].Description
			}

		})

		return result, nil

	} else {

		tool.mutex.Unlock()

		return nil, fmt.Errorf("changelog.Snapshot: %s", err.Error())

	}

}
