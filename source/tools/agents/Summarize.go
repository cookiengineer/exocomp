package agents

import "exocomp/schemas"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Agents) Summarize(messages []*schemas.Message) (string, error) {

	prompt := strings.Join([]string{
		"Please create a detailed technical summary of the following conversation.",
		"Preserve the user's goals, decisions, important file paths, symbols, tool results and unfinished work.",
		"The latest messages are the newest ones.",
		"",
		utils_chat.SummarizeMessages(messages, true, true, true),
	}, "\n")

	summary, err0 := tool.summarizeMessages(prompt)

	if err0 == nil {
		return summary, nil
	} else {
		return "", fmt.Errorf("agents.Switch: Failed to summarize previous conversation.")
	}

}
