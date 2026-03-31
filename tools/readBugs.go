package tools

import "exocomp/utils"
import "fmt"
import "os"
import "strings"

func readBugs(tool *Bugs) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./BUGS.md")

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				for anchor, _ := range tool.contents {
					delete(tool.contents, anchor)
				}

				lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")

				for _, line := range lines {

					if strings.TrimSpace(line) == "# Bugs" {

						continue

					} else if strings.TrimSpace(line) == "" {

						continue

					} else if strings.HasPrefix(line, "- [") && strings.Contains(line, "] `") && strings.Contains(line, "`: ") {

						tmp1 := strings.TrimSpace(line[2:strings.Index(line, "] `")])
						tmp2 := strings.TrimSpace(line[strings.Index(line, "] `"):strings.Index(line, "`: ")])
						tmp3 := strings.TrimSpace(line[strings.Index(line, "`: "):])

						anchor   := tmp2
						file     := ""
						note     := utils.FormatSingleLine(tmp3)
						// method   := ""
						is_fixed := false

						if tmp1 == "x" {
							is_fixed = true
						}

						if strings.Contains(tmp2, "#") {
							file   = strings.TrimSpace(tmp2[0:strings.Index(tmp2, "#")])
							// method = strings.TrimSpace(tmp2[strings.Index(tmp2, "#"):])
						} else {
							file   = strings.TrimSpace(tmp2)
							// method = ""
						}

						if file != "" && note != "" {

							_, err2 := resolveSandboxPath(tool.Sandbox, file)

							if err2 == nil {

								_, ok1 := tool.contents[anchor]

								if ok1 == false {
									tool.contents[anchor] = make(map[string]bool)
								}

								tool.contents[anchor][note] = is_fixed

							}

						}

					}

				}

				return nil

			} else {
				return fmt.Errorf("readBugs: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readBugs: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readBugs: Invalid Tool Sandbox")
	}

}
