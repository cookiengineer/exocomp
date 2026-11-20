package bugs

import "exocomp/types"
import "fmt"
import "sort"

func (tool *Bugs) Snapshot() ([]types.Bug, error) {

	tool.mutex.Lock()
	err := readBugs(tool)

	if err == nil {

		result := make([]types.Bug, 0)

		for _, bug_reports := range tool.contents {

			for _, bug_report := range bug_reports {

				if bug_report.IsFixed == false {
					result = append(result, bug_report)
				}

			}

		}

		tool.mutex.Unlock()

		sort.Slice(result, func(a int, b int) bool {

			if result[a].File != result[b].File {
				return result[a].File < result[b].File
			} else if result[a].Symbol != result[b].Symbol {
				return result[a].Symbol < result[b].Symbol
			} else {
				return result[a].Description < result[b].Description
			}

		})

		return result, nil

	} else {

		tool.mutex.Unlock()

		return nil, fmt.Errorf("bugs.OpenBugs: %s", err.Error())

	}

}
