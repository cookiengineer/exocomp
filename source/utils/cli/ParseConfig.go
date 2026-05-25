package cli

import "exocomp/types"
import utils_agents "exocomp/utils/agents"
import utils_fmt "exocomp/utils/fmt"
import net_url "net/url"
import "os"
import "strconv"
import "strings"

func ParseConfig(arguments []string) *types.Config {

	name          := ""
	role          := "planner"
	debug         := false
	model         := "qwen3-coder:30b"
	playground, _ := os.Getwd()
	prompt        := ""
	sandbox, _    := os.Getwd()
	temperature   := float64(0.0)
	url, _        := net_url.Parse("http://localhost:11434/v1")

	for _, argument := range arguments {

		if strings.HasPrefix(argument, "--") && strings.Contains(argument, "=") {

			flag := strings.Split(argument[2:], "=")

			if len(flag) == 2 {

				if strings.HasPrefix(flag[1], "\"") && strings.HasSuffix(flag[1], "\"") {
					flag[1] = flag[1][1:len(flag[1]) - 1]
				}

				if strings.HasPrefix(flag[1], "'") && strings.HasSuffix(flag[1], "'") {
					flag[1] = flag[1][1:len(flag[1]) - 1]
				}

				switch flag[0] {
				case "name":

					if utils_agents.IsAgentName(flag[1]) {
						name = utils_fmt.FormatAgentName(flag[1])
					}

				case "role":

					if utils_agents.IsAgentRole(flag[1]) {
						role = utils_fmt.FormatAgentRole(flag[1])
					}

				case "model":

					model = strings.TrimSpace(flag[1])

				case "playground":

					stat, err := os.Stat(strings.TrimSpace(flag[1]))

					if err == nil && stat.IsDir() {
						playground = strings.TrimSpace(flag[1])
					} else if err != nil && os.IsNotExist(err) {
						playground = strings.TrimSpace(flag[1])
					}

				case "prompt":

					prompt = utils_fmt.FormatMultiLine(flag[1])

				case "sandbox":

					stat, err := os.Stat(strings.TrimSpace(flag[1]))

					if err == nil && stat.IsDir() {
						sandbox = strings.TrimSpace(flag[1])
					} else if err != nil && os.IsNotExist(err) {
						sandbox = strings.TrimSpace(flag[1])
					}

				case "temperature":

					tmp, err := strconv.ParseFloat(strings.TrimSpace(flag[1]), 10)

					if err == nil {

						if tmp >= 0.1 && tmp <= 1.0 {
							temperature = tmp
						}

					}

				case "url":

					tmp, err := net_url.Parse(strings.TrimSpace(flag[1]))

					if err == nil {

						if tmp.Scheme == "http" || tmp.Scheme == "https" {

							if tmp.Path == "/" || tmp.Path == "/v1" {
								url = tmp
							}

						}

					}

				}

			}

		} else if strings.HasPrefix(argument, "--") {

			flag := strings.TrimSpace(argument[2:])

			switch flag {
			case "debug":
				debug = true
			}

		}

	}

	if playground == sandbox || strings.HasPrefix(sandbox, playground + string(os.PathSeparator)) {

		return types.NewConfig(
			name,
			role,
			model,
			prompt,
			temperature,
			playground,
			sandbox,
			url,
			debug,
		)

	} else {
		return nil
	}

}
