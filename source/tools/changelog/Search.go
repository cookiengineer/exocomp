package changelog

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"
import "time"

func (tool *Changelog) Search(path string, symbol string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readChangelog(tool)

			found := make(map[time.Time][]string, 0)

			if symbol != "" {

				tool.mutex.RLock()

				_, ok1 := tool.contents[internal_path]

				if ok1 == true {

					entries, ok2 := tool.contents[internal_path][symbol]

					if ok2 == true {

						for _, entry := range entries {

							sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, entry.File)

							if err3 == nil {
								found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: \"%s\", File: \"%s\", Symbol: \"%s\", Description: \"%s\"", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
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
				result = append(result, fmt.Sprintf("changelog.Search: File \"%s\" with Symbol \"%s\" contains %d Changelog entries.", path, symbol, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			} else {

				tool.mutex.RLock()

				symbols, ok1 := tool.contents[internal_path]

				if ok1 == true {

					sorted_symbols := make([]string, 0)

					for symbol, _ := range symbols {
						sorted_symbols = append(sorted_symbols, symbol)
					}

					sort.Strings(sorted_symbols)

					for _, symbol := range sorted_symbols {

						entries := tool.contents[internal_path][symbol]

						for _, entry := range entries {

							sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, entry.File)

							if err3 == nil {
								found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: %s, File: %s, Symbol: %s, Description: %s", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
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
				result = append(result, fmt.Sprintf("changelog.Search: File \"%s\" contains %d Changelog entries.", path, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			}

		} else {
			return "", fmt.Errorf("changelog.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.Search: %s", err1.Error())
	}

}
