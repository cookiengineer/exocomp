package tools

import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "strings"
import "time"

func handleProcessJSONLine(tool *Agents, name string, line []byte) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	process, ok := tool.processes[name]

	if ok == true {
		process.last_write = time.Now()
	}

	text := strings.TrimRight(string(line), "\n")

	if strings.HasPrefix(text, "schemas.Message:") {

		buffer := []byte(text[16:])
		message := schemas.Message{}

		err := json.Unmarshal(buffer, &message)

		if err == nil {

			agent, ok1 := tool.contents[name]

			if ok1 == true {
				agent.Messages = append(agent.Messages, &message)
			}

		}

	} else if strings.HasPrefix(text, "types.ContextUsage:") {

		buffer := []byte(text[19:])
		usage := types.ContextUsage{}

		err := json.Unmarshal(buffer, &usage)

		if err == nil {

			agent, ok1 := tool.contents[name]

			if ok1 == true {

				agent.ContextUsage.Length = usage.Length
				agent.ContextUsage.Tokens = usage.Tokens

			}

		}

	}

}
