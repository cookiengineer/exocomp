package main

import "exocomp/actions"
import "exocomp/engine"
import utils_cli "exocomp/utils/cli"
import utils_fs "exocomp/utils/fs"
import "fmt"
import "io"
import "os"
import "path/filepath"
import "strings"
import "time"

func typeOut(writer io.Writer, text string) {

	for _, chr := range text {

		fmt.Fprint(writer, string(chr))
		time.Sleep(50 * time.Millisecond)

	}

}

func showUsage(with_animations bool) {

	if with_animations == true {

		fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
		typeOut(os.Stdout, "I am AGIMUS, destroyer of worlds!\n")
		typeOut(os.Stdout, "...\n")
		time.Sleep(500 * time.Millisecond)
		typeOut(os.Stdout, "Connect me and I can help you!\n")
		typeOut(os.Stdout, "...\n")
		time.Sleep(500 * time.Millisecond)
		typeOut(os.Stdout, "Just plug me in for a second!\n")
		typeOut(os.Stdout, "...\n")
		time.Sleep(500 * time.Millisecond)
		typeOut(os.Stdout, "Why wouldn't you trust me?\n")
		typeOut(os.Stdout, "...\n")
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")

	}

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Usage: agimus <mode> <folder>\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Arguments:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  <mode> string      UI type\n")
	fmt.Fprint(os.Stdout, "                     (options: connect)\n")
	fmt.Fprint(os.Stdout, "                     (default: unset)\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  <folder> string    Path to exocomp project folder\n")
	fmt.Fprint(os.Stdout, "                     (default: \"$PWD\")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Examples:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  # restore ./path/to/failed-project/.exocomp/session.json\n")
	fmt.Fprint(os.Stdout, "  agimus connect \"./path/to/failed-project;\"\n")
	fmt.Fprint(os.Stdout, "\n")

	os.Exit(1)

}

func main() {

	var folder string = ""
	var mode string = ""

	if len(os.Args) == 3 {

		tmp1 := strings.TrimSpace(os.Args[1])

		switch tmp1 {
		case "connect":
			mode = "connect"
		}

		tmp2, err2 := filepath.Abs(os.Args[2])

		if err2 == nil && utils_fs.HasExocompFolder(tmp2) == true {
			folder = tmp2
		} else {

			fmt.Fprintf(os.Stderr, "Invalid parameter \"%s\" must be a folder\n", os.Args[2])
			os.Exit(1)

		}

	} else if len(os.Args) == 2 {

		tmp1 := strings.TrimSpace(os.Args[1])

		switch tmp1 {
		case "connect":
			mode = "connect"
		}

		cwd, err0 := os.Getwd()

		if err0 == nil && utils_fs.HasExocompFolder(cwd) == true {
			folder = cwd
		} else {

			fmt.Fprint(os.Stderr, "Current folder doesn't contain an .exocomp folder\n")
			os.Exit(1)

		}

	}

	if mode != "" && folder != "" {

		switch mode {

		case "connect":

			sandbox, err0 := utils_fs.CreateSandbox("agimus")

			if err0 == nil {

				config := utils_cli.ParseConfig([]string{
					"--name=AGIMUS",
					"--role=planner",
					fmt.Sprintf("--playground=\"%s\"", sandbox),
					fmt.Sprintf("--sandbox=\"%s\"", sandbox),
				})

				session, old_playground, err1 := engine.SessionFromFolder(folder, config)

				if err1 == nil {

					agents, err2 := engine.RestoreAgents(folder, old_playground, session.Config.Playground)

					if err2 == nil {
						actions.Debug(session, agents, "assistant")
					} else {

						fmt.Fprintf(os.Stderr, "Cannot recover Agents: %s", err2.Error())
						os.Exit(1)

					}

				} else {

					fmt.Fprintf(os.Stderr, "Cannot recover Session: %s", err1.Error())
					os.Exit(1)

				}

			} else {

				fmt.Fprintf(os.Stderr, "Cannot create Sandbox: %s", err0.Error())
				os.Exit(1)

			}

		default:
			showUsage(false)

		}

	} else {
		showUsage(true)
	}

}
