package actions

import "exocomp/agents"
import "exocomp/types"
import "fmt"
import "os"
import "sort"
import "strconv"
import "strings"

func Usage(agent *types.Agent, config *types.Config, options_ui []string) {

	default_ui          := "unset"
	default_agent       := strconv.Quote("planner")
	default_name        := strconv.Quote("Peanut Hamper")
	default_model       := strconv.Quote("qwen3-coder:30b")
	default_prompt      := "unset"
	default_temperature := "unset"
	default_sandbox     := "current working directory"
	default_url         := strconv.Quote("http://localhost:11434/v1")

	options_agent := make([]string, 0)
	options_model := []string{
		strconv.Quote("gemma4:31b"),
		strconv.Quote("qwen3-coder:30b"),
		strconv.Quote("qwen3.6:35b-heretic"),
	}
	options_temperature := []string{
		strconv.FormatFloat(0.1, 'f', 1, 64),
		strconv.FormatFloat(1.0, 'f', 1, 64),
	}

	for _, agent := range agents.Types {
		options_agent = append(options_agent, strconv.Quote(agent))
	}

	sort.Strings(options_agent)
	sort.Strings(options_model)

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Usage: exocomp <ui> [flags]\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Arguments:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    <ui> string            UI type\n")
	fmt.Fprint(os.Stdout, "                           (options: " + strings.Join(options_ui, ", ") + ")\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_ui + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Flags:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --name string          LLM agent name\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_name + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --agent string         LLM agent type\n")
	fmt.Fprint(os.Stdout, "                           (options: " + strings.Join(options_agent, ", ") + ")\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_agent + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --model string         LLM agent model\n")
	fmt.Fprint(os.Stdout, "                           (options: " + strings.Join(options_model, ", ") + ")\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_model + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --temperature float    LLM agent sampling temperature\n")
	fmt.Fprint(os.Stdout, "                           (options: " + strings.Join(options_temperature, "-") + ")\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_temperature + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --prompt string        Initial LLM instructions prompt\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_prompt + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --sandbox string       Path to sandbox directory\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_sandbox + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    --url string           OpenAI API endpoint for LLM backend\n")
	fmt.Fprint(os.Stdout, "                           (default: " + default_url + ")\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Examples:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    # single-agent mode\n")
	fmt.Fprint(os.Stdout, "    exocomp terminal --agent=architect;\n")
	fmt.Fprint(os.Stdout, "    exocomp web --agent=architect --model=\"qwen3-coder:30b\" --temperature=\"0.5\";\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "    # multi-agent mode\n")
	fmt.Fprint(os.Stdout, "    exocomp terminal --agent=planner;\n")
	fmt.Fprint(os.Stdout, "    exocomp web --agent=planner --model=\"codestral:22b\" --temperature=\"0.7\";\n")
	fmt.Fprint(os.Stdout, "\n")

	os.Exit(1)

}
