package tools

import "time"

func watchProcess(tool *Agents, process *agentProcess) {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-process.done:
			return

		case <-ticker.C:

			tool.mutex.Lock()
			stale := tool.IdleTimeout > 0 && time.Since(process.last_write) > tool.IdleTimeout
			tool.mutex.Unlock()

			if stale == true {

				process.cancel()
				return

			}

		}

	}

}
