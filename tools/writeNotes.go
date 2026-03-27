package tools

import "fmt"
import "os"
import "slices"
import "strings"

func writeNotes(tool *Notes) error {

	if tool.Sandbox != "" {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, "./NOTES.md")

		if err0 == nil {

			ids := make([]uint, 0)

			for id, _ := range tool.notes {
				ids = append(ids, id)
			}

			slices.Sort(ids)

			content := make([]string, 0)

			content = append(content, "# Notes")
			content = append(content, "\n")

			for _, id := range ids {
				content = append(content, fmt.Sprintf("- %s\n", tool.notes[id]))
			}

			content = append(content, "\n")

			bytes := []byte(strings.Join(content, "\n"))
			err1  := os.WriteFile(resolved, bytes, 0666)

			if err1 == nil {
				return nil
			} else {
				return fmt.Errorf("writeNotes: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("writeNotes: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("writeNotes: Invalid Tool Sandbox")
	}

}
