package cli

import "exocomp/types"
import utils_agent "exocomp/utils/agent"
import utils_fmt "exocomp/utils/fmt"
import net_url "net/url"
import "os"
import "strconv"
import "strings"

func ParseConfig(arguments []string) *types.Config {

	name          := ""
	agent         := "planner"
	debug         := false
	model         := "qwen3-coder:30b"
	playground, _ := os.Getwd()
	prompt        := ""
	sandbox, _    := os.Getwd()
	temperature   := float64(0.0)
	url, _        := net_url.Parse("http://localhost:11434/v1")

	for _, argument := range arguments {

		if strings.HasPrefix(argument, "--") && strings.Contains(argument, "=") {

			tmp := strings.Split(argument[2:], "=")

			if len(tmp) == 2 {

				if strings.HasPrefix(tmp[1], "\"") && strings.HasSuffix(tmp[1], "\"") {
					tmp[1] = tmp[1][1:len(tmp[1]) - 2]
				}

				if strings.HasPrefix(tmp[1], "'") && strings.HasSuffix(tmp[1], "'") {
					tmp[1] = tmp[1][1:len(tmp[1]) - 2]
				}

				switch tmp[0] {
				case "name":

					if utils_agent.IsName(tmp[1]) {
						name = utils_fmt.FormatAgentName(tmp[1])
					}

				case "agent":

					if utils_agent.IsType(tmp[1]) {
						agent = strings.TrimSpace(tmp[1])
					}

				case "model":

					model = strings.TrimSpace(tmp[1])

				case "playground":

					stat, err := os.Stat(strings.TrimSpace(tmp[1]))

					if err == nil && stat.IsDir() {
						playground = strings.TrimSpace(tmp[1])
					} else if err != nil && os.IsNotExist(err) {
						playground = strings.TrimSpace(tmp[1])
					}

				case "prompt":

					prompt = utils_fmt.FormatSingleLine(tmp[1])

				case "sandbox":

					stat, err := os.Stat(strings.TrimSpace(tmp[1]))

					if err == nil && stat.IsDir() {
						sandbox = strings.TrimSpace(tmp[1])
					} else if err != nil && os.IsNotExist(err) {
						sandbox = strings.TrimSpace(tmp[1])
					}

				case "temperature":

					num, err := strconv.ParseFloat(strings.TrimSpace(tmp[1]), 10)

					if err == nil {

						if num >= 0.1 && num <= 1.0 {
							temperature = num
						}

					}

				case "url":

					url, err := net_url.Parse(strings.TrimSpace(tmp[1]))

					if err == nil {

						if url.Scheme == "http" || url.Scheme == "https" {

							if url.Path == "/" || url.Path == "/v1" {
								url = url
							}

						}

					}

				}

			}

		} else if strings.HasPrefix(argument, "--") {

			tmp := strings.TrimSpace(argument[2:])

			switch tmp {
			case "debug":
				debug = true
			}

		}

	}

	if playground == sandbox || strings.HasPrefix(sandbox, playground + string(os.PathSeparator)) {

		return types.NewConfig(
			name,
			agent,
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
