package tools

import "fmt"
import "os"
import "strings"

func writeBugs(tool *Bugs) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./BUGS.md")

		if err0 == nil {

			lines := make([]string, 0)

			lines = append(lines, "# Bugs")
			lines = append(lines, "")

			for anchor, notes := range tool.contents {

				for note, is_fixed := range notes {

					if is_fixed == true {
						lines = append(lines, fmt.Sprintf("- [%s] `%s`: %s\n", "x", anchor, note))
					} else {
						lines = append(lines, fmt.Sprintf("- [%s] `%s`: %s\n", " ", anchor, note))
					}

				}

			}

			bytes := []byte(strings.Join(lines, "\n") + "\n")
			err1  := os.WriteFile(resolved, bytes, 0666)

			if err1 == nil {
				return nil
			} else {
				return fmt.Errorf("writeBugs: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("writeBugs: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("writeBugs: Invalid Tool Sandbox")
	}

}
