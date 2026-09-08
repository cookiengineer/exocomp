package engine

import "exocomp/types"
import "fmt"
import "os"
import "path/filepath"
import "strings"

func RestoreAgents(folder string, old_playground string, new_playground string) ([]*types.Agent, error) {

	result := make([]*types.Agent, 0)

	entries, err1 := os.ReadDir(filepath.Join(folder, ".exocomp", "agents"))

	if err1 == nil {

		errors := make([]string, 0)

		for _, entry := range entries {

			filename := entry.Name()

			if strings.HasSuffix(filename, ".json") {

				agentname := strings.TrimSpace(filename[0:len(filename)-5])

				if agentname != "" {

					bytes, err2 := os.ReadFile(filepath.Join(folder, ".exocomp", "agents", filename))

					if err2 == nil {

						agent, err3 := types.ParseAgent(bytes)

						if err3 == nil {

							if old_playground != "" && new_playground != "" && agent.Sandbox != "" {

								relative, err4 := filepath.Rel(old_playground, agent.Sandbox)

								if err4 == nil && relative != ".." && strings.HasPrefix(relative, ".."+string(os.PathSeparator)) == false {
									agent.Sandbox = filepath.Join(new_playground, relative)
								}

							}

							result = append(result, agent)

						} else {
							errors = append(errors, fmt.Sprintf("%s: %s", filename, err3.Error()))
						}

					} else {
						errors = append(errors, fmt.Sprintf("%s: %s", filename, err2.Error()))
					}

				}

			}

		}

		if len(errors) == 0 {
			return result, nil
		} else {
			return result, fmt.Errorf("%s", strings.Join(errors, "\n"))
		}

	} else {
		return result, err1
	}

}
