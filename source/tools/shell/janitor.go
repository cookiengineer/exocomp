package shell

import "time"

func (tool *Shell) janitor() {

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {

		select {
		case <-tool.stop:
			return
		case <-ticker.C:
			tool.cleanup()
		}

	}

}
