package agents

import "exocomp/schemas"
import utils_cli "exocomp/utils/cli"
import utils_rand "exocomp/utils/rand"
import "bytes"
import "context"
import "encoding/json"
import "fmt"
import "os"
import "os/exec"
import "strings"
import "time"

func (tool *Agents) summarizeMessages(prompt string) (string, error) {

	tmp, err0 := os.MkdirTemp("/tmp", "exocomp-summarizer-*")

	if err0 == nil {

		exe, err1 := utils_cli.FindAgentBinary()

		if err1 == nil {

			timeout := tool.Timeout

			if timeout <= 0 {
				timeout = 30 * time.Minute
			}

			resolved_url := ""

			if tool.URL != nil {
				resolved_url = tool.URL.String()
			}

			arguments := []string{
				fmt.Sprintf("--name=\"%s\"", utils_rand.AgentName("summarizer")),
				fmt.Sprintf("--role=\"%s\"", "summarizer"),
				fmt.Sprintf("--model=\"%s\"", tool.Model),
				fmt.Sprintf("--prompt=\"%s\"", prompt),
				// --temperature set by agent role
				// --playground set by cmd.Dir
				// --sandbox set by cmd.Dir
				fmt.Sprintf("--url=\"%s\"", resolved_url),
			}

			if tool.Debug == true {
				arguments = append(arguments, "--debug")
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			cmd := exec.CommandContext(ctx, exe, arguments...)
			cmd.Dir = tmp

			stdout_buffer := bytes.Buffer{}
			cmd.Stdout = &stdout_buffer

			err2 := cmd.Run()

			if err2 == nil {

				os.RemoveAll(tmp)

				lines := strings.Split(strings.TrimSpace(stdout_buffer.String()), "\n")

				if len(lines) > 0 {

					line := strings.TrimSpace(lines[len(lines)-1])
					summary := schemas.Message{}

					if strings.HasPrefix(line, "schemas.Message:") {
						line = strings.TrimSpace(line[16:])
					}

					err3 := json.Unmarshal([]byte(line), &summary)

					if err3 == nil {
						return strings.TrimSpace(summary.Content), nil
					} else {
						return "", fmt.Errorf("agents.Summarize: Cannot parse summary.")
					}

				} else {
					return "", fmt.Errorf("agents.Summarize: Cannot parse summary.")
				}

			} else {

				os.RemoveAll(tmp)
				return "", fmt.Errorf("agents.Summarize: Cannot summarize messages.")

			}

		} else {

			os.RemoveAll(tmp)
			return "", fmt.Errorf("agents.Summarize: exocomp-agent is not installed.")

		}

	} else {
		return "", fmt.Errorf("agents.Summarize: Cannot create temporary workspace.")
	}

}
