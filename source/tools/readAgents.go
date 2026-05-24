package tools

import "exocomp/agents"
import "exocomp/types"
import "fmt"
import "os"
import "path/filepath"
import "strings"

func readAgents(tool *Agents) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, "agents")

		if err0 == nil {

			agents_entries, err1 := os.ReadDir(resolved)

			if err1 == nil {

				errors := make([]error, 0)

				for _, agent_entry := range agents_entries {

					agent_filename := agent_entry.Name()
					ext            := filepath.Ext(agent_filename)
					agent_name     := strings.TrimSuffix(agent_filename, ext)

					if ext == ".yaml" {

						agent_path       := filepath.Join(tool.Playground, "agents", agent_filename)
						agent_stat, err1 := os.Stat(agent_path)

						if err1 == nil && agent_stat.IsDir() == false {

							agent_bytes, err2 := os.ReadFile(agent_path)

							if err2 == nil {

								agent, err3 := types.ParseAgent(agent_bytes)

								if err3 == nil {

									if agent.Role != "" {
										agents.SetRole(agent.Role, agent)
									}

								} else {
									errors = append(errors, fmt.Errorf("Invalid Agent %s: Cannot parse %s", agent_name, agent_filename))
								}

							} else {
								errors = append(errors, fmt.Errorf("Invalid Agent %s: Cannot parse %s", agent_name, agent_filename))
							}

						}

					}

				}

				if len(errors) > 0 {

					error_messages := make([]string, 0)

					for _, err := range errors {
						error_messages = append(error_messages, strings.TrimSpace(err.Error()))
					}

					return fmt.Errorf("readAgents: %s", strings.Join(error_messages, " | "))

				} else {
					return nil
				}

			} else {
				return fmt.Errorf("readAgents: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readAgents: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readAgents: Invalid tool playground")
	}

}
