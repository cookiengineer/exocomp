package changelog

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"
import "time"

func (tool *Changelog) List() (string, error) {

	readChangelog(tool)

	found := make(map[time.Time][]string, 0)

	tool.mutex.RLock()

	for _, symbols := range tool.contents {

		for _, entries := range symbols {

			for _, entry := range entries {

				resolved_path, err1 := utils_sandbox.ResolvePath(tool.Playground, entry.File)

				if err1 == nil {

					sandbox_path, err2 := utils_sandbox.SanitizePath(tool.Sandbox, resolved_path)

					if err2 == nil {
						found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: \"%s\", File: \"%s\", Symbol: \"%s\", Description: \"%s\"", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
					}

				}

			}

		}

	}

	tool.mutex.RUnlock()

	dates := make([]time.Time, 0)
	lines := make([]string, 0)

	for date, _ := range found {
		dates = append(dates, date)
	}

	sort.Slice(dates, func(a int, b int) bool {
		return dates[a].Before(dates[b])
	})

	for _, date := range dates {

		for _, line := range found[date] {
			lines = append(lines, line)
		}

	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("changelog.List: %d Changelog entries.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}
