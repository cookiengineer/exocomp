package agents

import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Agents) Inquire(name string) (string, error) {

	tool.mutex.Lock()
	agent, ok0 := tool.contents[name]
	tool.mutex.Unlock()

	if ok0 == true {

		tool.mutex.Lock()
		messages := utils_chat.SummarizeMessages(agent.Messages, true, true, false)
		tool.mutex.Unlock()

		prompt := strings.Join([]string{
			"Please summarize the following conversation, the latest messages are the newest ones.",
			"",
			messages,
		}, "\n")

		summary, err0 := tool.summarizeMessages(prompt)

		if err0 == nil {

			tool.mutex.Lock()
			_, ok1 := tool.processes[name]
			tool.mutex.Unlock()

			if ok1 == true {

				result := strings.Join([]string{
					fmt.Sprintf("agents.Inquire: Summary of currently working agent \"%s\"'s work report:", name),
					strings.TrimSpace(summary),
				}, "\n")

				return result, nil

			} else {

				result := strings.Join([]string{
					fmt.Sprintf("agents.Inquire: Summary of already finished agent \"%s\"'s work report:", name),
					strings.TrimSpace(summary),
				}, "\n")

				return result, nil

			}

		} else {
			return "", fmt.Errorf("agents.Inquire: Failed to summarize agent \"%s\"'s work report!", name)
		}

	} else {
		return "", fmt.Errorf("agents.Inquire: Agent \"%s\" didn't work for us!", name)
	}

}
