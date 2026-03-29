package tools

import "fmt"
import "os"
import "strings"
import "time"

func readChangelog(tool *Changelog) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./CHANGELOG.md")

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				for date, _ := range tool.contents {
					delete(tool.contents, date)
				}

				lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
				today := time.Time{}

				for _, line := range lines {

					if strings.TrimSpace(line) == "# Changelog" {

						continue

					} else if strings.TrimSpace(line) == "" {

						continue

					} else if strings.HasPrefix(line, "## ") {

						tmp1      := strings.TrimSpace(line[3:])
						tmp2, err := time.Parse("2006-01-02 15:04:05", tmp1)

						if err == nil {
							today = tmp2
						}

					} else if strings.HasPrefix(line, "`") && strings.Contains(line, "`: ") {

						if today.IsZero() == false {

							_, err := resolveSandboxPath(tool.Sandbox, line[1:strings.Index(line, "`:")])

							if err == nil {

								_, ok := tool.contents[today]

								if ok == false {
									tool.contents[today] = make([]string, 0)
								}

								tool.contents[today] = append(tool.contents[today], strings.TrimSpace(line))

							}

						}

					}

				}

				return nil

			} else {
				return fmt.Errorf("readChangelog: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readChangelog: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readChangelog: Invalid Tool Sandbox")
	}

}
