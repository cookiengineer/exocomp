package tools

import "fmt"
import "os"
import "strings"

func readNotes(tool *Notes) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./NOTES.md")

		if err0 == nil {

			bytes, err1 := os.ReadFile(resolved)

			if err1 == nil {

				lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
				notes := make([]string, 0)

				for _, line := range lines {

					if strings.TrimSpace(line) == "# Notes" {

						continue

					} else if strings.TrimSpace(line) == "" {

						continue

					} else if strings.HasPrefix(line, "- ") {

						note := strings.TrimSpace(line[2:])

						if note != "" {
							notes = append(notes, note)
						}

					}

				}

				for old_id, _ := range tool.notes {
					delete(tool.notes, old_id)
				}

				for n, note := range notes {
					tool.notes[uint(n)] = note
				}

				return nil

			} else {
				return fmt.Errorf("readNotes: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readNotes: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readNotes: Invalid Tool Sandbox")
	}

}
