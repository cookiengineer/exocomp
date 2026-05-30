package main

import "exocomp/actions"
import "exocomp/agents"
import "exocomp/tools"
import "exocomp/types"
import utils_cli "exocomp/utils/cli"
import ui_web "exocomp/ui/web"
import "fmt"
import "io"
import "os"
import "path/filepath"
import "slices"
import "strings"
import "time"

func restore_agents(folder string) ([]*types.Agent, error) {

	result := make([]*types.Agent, 0)

	entries, err1 := os.ReadDir(filepath.Join(folder, "agents"))

	if err1 == nil {

		errors := make([]string, 0)

		for _, entry := range entries {

			filename := entry.Name()

			if strings.HasSuffix(filename, ".json") {

				agentname := strings.TrimSpace(filename[0:len(filename)-5])

				if agentname != "" {

					bytes, err2 := os.ReadFile(filepath.Join(folder, "agents", filename))

					if err2 == nil {

						agent, err3 := types.ParseAgent(bytes)

						if err3 == nil {
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

func restore_session(folder string) (*types.Session, error) {

	bytes, err1 := os.ReadFile(filepath.Join(folder, "session.json"))

	if err1 == nil {

		backup, err2 := types.ParseSession(bytes)

		if err2 == nil {

			cwd, err3 := os.Getwd()

			if err3 == nil {
				return types.RestoreSession(cwd, *backup), nil
			} else {
				return nil, err3
			}

		} else {
			return nil, err2
		}

	} else {
		return nil, err1
	}

}

func type_out(writer io.Writer, text string) {

	for _, chr := range text {
		fmt.Fprint(writer, string(chr))
		time.Sleep(50 * time.Millisecond)
	}

}

func showUsage() {

	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
	type_out(os.Stdout, "I am AGIMUS, destroyer of worlds!\n")
	type_out(os.Stdout, "...\n")
	time.Sleep(1500 * time.Millisecond)
	type_out(os.Stdout, "Connect me and I can help you!\n")
	type_out(os.Stdout, "...\n")
	time.Sleep(1000 * time.Millisecond)
	type_out(os.Stdout, "Just plug me in for a second!\n")
	type_out(os.Stdout, "...\n")
	time.Sleep(1000 * time.Millisecond)
	type_out(os.Stdout, "Why wouldn't you trust me?\n")
	type_out(os.Stdout, "...\n")
	time.Sleep(1000 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Usage: agimus <action> <folder>\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Actions:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  connect    tests a project's sandboxes\n")
	fmt.Fprint(os.Stdout, "             (default: \"$PWD/.exocomp\")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  inspect    inspects a project's session\n")
	fmt.Fprint(os.Stdout, "             (default: \"$PWD/.exocomp\")\n")
	fmt.Fprint(os.Stdout, "\n")

}

func main() {

	options := []string{"connect", "inspect"}
	action  := ""
	folder  := ""

	if len(os.Args) == 3 {

		if slices.Contains(options, os.Args[1]) {
			action = os.Args[1]
		}

		if strings.HasSuffix(os.Args[2], ".exocomp") {

			stat, err0 := os.Stat(os.Args[2])

			if err0 == nil && stat.IsDir() {

				tmp, err1 := filepath.Abs(os.Args[2])

				if err1 == nil {
					folder = tmp
				} else {

					fmt.Fprintf(os.Stderr, "Error: %s", err1)
					fmt.Fprintf(os.Stderr, "Invalid parameter \"%s\" must be an .exocomp folder\n", os.Args[2])
					os.Exit(1)

				}

			} else {

				fmt.Fprintf(os.Stderr, "Error: %s", err0)
				fmt.Fprintf(os.Stderr, "Invalid parameter \"%s\" must be an .exocomp folder\n", os.Args[2])
				os.Exit(1)

			}

		} else {

			fmt.Fprintf(os.Stderr, "Invalid parameter \"%s\" must be an .exocomp folder\n", os.Args[2])
			os.Exit(1)

		}

	} else if len(os.Args) == 2 {

		if slices.Contains(options, os.Args[1]) {
			action = os.Args[1]
		}

		cwd, err0 := os.Getwd()

		if err0 == nil {

			tmp        := filepath.Join(cwd, ".exocomp")
			stat, err1 := os.Stat(tmp)

			if err1 == nil && stat.IsDir() {
				folder = tmp
			} else {

				fmt.Fprint(os.Stderr, "Current folder doesn't contain an .exocomp folder\n")
				os.Exit(1)

			}

		}

	}

	if action != "" && folder != "" {

		if action == "connect" {

			config := utils_cli.ParseConfig([]string{"planner"})
			config.Name = "AGIMUS"

			agent := agents.NewAgent(config)

			config.Update(
				agent.Name,
				agent.Role,
				agent.Model,
				config.Prompt,
				agent.Temperature,
			)

			err0 := os.MkdirAll(config.Sandbox, 0755)

			if err0 == nil {

				actions.Terminal(agent, config, "assistant")

			} else {

				fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
				os.Exit(1)

			}

		} else if action == "inspect" {

			session, err1 := restore_session(folder)

			if err1 == nil {

				agents, err2 := restore_agents(folder)

				if err2 != nil {
					fmt.Fprintf(os.Stderr, "Error: %s", err2.Error())
				}

				server := ui_web.NewServer(session.Agent, session.Config)

				// Override default recovery
				server.Session = session

				if len(session.Agent.AllowedTools) > 0 {

					tool_schemas, tools := tools.Toolset(
						session.Config.Playground,
						session.Config.Sandbox,
						session.Config.Model,
						session.Config.URL,
						session.Config.Debug,
						session.Agent.AllowedPrograms,
						session.Agent.AllowedTools,
					)

					for name, tool := range tools {
						server.Session.SetTool(name, tool, tool_schemas[name])
					}

				}

				go func() {

					time.Sleep(100 * time.Millisecond)

					tool := server.Session.GetTool("agents.List")

					if tool != nil {

						agent_tool, ok := tool.(*tools.Agents)

						if ok == true {

							fmt.Fprintf(os.Stdout, "|\n")

							for _, agent := range agents {
								fmt.Fprintf(os.Stdout, "|-> Restored Agent \"%s\" with %d messages\n", agent.Name, len(agent.Messages))
								agent_tool.SetAgent(agent)
							}

						}

					}

				}()

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", session.Agent.Name, session.Agent.Role, session.Agent.Model, session.Agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", session.Config.Sandbox)
				fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(session.Agent.AllowedTools, ", "))
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", session.Config.URL.String())
				fmt.Fprintf(os.Stdout, "| Web:     %s\n", server.URL.String())

				server.Init()

			} else {
				fmt.Fprintf(os.Stderr, "Error: %s", err1.Error())
				os.Exit(1)
			}

		}

	} else {

		showUsage()
		os.Exit(1)

	}

}
