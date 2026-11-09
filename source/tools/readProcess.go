package tools

import "exocomp/schemas"
import "bufio"
import "io"

func readProcess(tool *Agents, name string, process *agentProcess, stdout_pipe io.ReadCloser) {

	reader := bufio.NewReader(stdout_pipe)

	for {

		line, err := reader.ReadBytes('\n')

		if len(line) > 0 {
			handleProcessJSONLine(tool, name, line)
		}

		if err != nil {
			break
		}

	}

	stdout_pipe.Close()

	process.cmd.Wait()

	tool.mutex.Lock()

	agent, ok := tool.contents[name]

	if ok == true && agent.Status == "working" {

		if getAgentWorkReport(agent) != "" {

			agent.FinishedAt = schemas.NewDatetime()
			agent.Status = "finished"

		} else {

			agent.FinishedAt = schemas.NewDatetime()
			agent.Status = "failed"

		}

	}

	delete(tool.processes, name)

	tool.mutex.Unlock()

	close(process.done)

}
