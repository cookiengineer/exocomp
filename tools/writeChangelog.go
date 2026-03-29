package tools

import "fmt"
import "os"
import "sort"
import "strings"
import "time"

func writeChangelog(tool *Changelog) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./CHANGELOG.md")

		if err0 == nil {

			lines := make([]string, 0)

			lines = append(lines, "# Changelog")
			lines = append(lines, "")

			dates := make([]time.Time, 0)

			for date, _ := range tool.contents {
				dates = append(dates, date)
			}

			sort.Slice(dates, func(a int, b int) bool {
				return dates[a].After(dates[b])
			})

			for _, date := range dates {

				entries := tool.contents[date]

				lines = append(lines, fmt.Sprintf("## %s", date.Format("2006-01-02 15:04:05")))
				lines = append(lines, "")

				sort.Strings(entries)

				for _, note := range entries {
					lines = append(lines, fmt.Sprintf("%s", note))
				}

				lines = append(lines, "")

			}

			bytes := []byte(strings.Join(lines, "\n") + "\n")
			err1  := os.WriteFile(resolved, bytes, 0666)

			if err1 == nil {
				return nil
			} else {
				return fmt.Errorf("writeChangelog: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("writeChangelog: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("writeChangelog: Invalid Tool Sandbox")
	}

}
