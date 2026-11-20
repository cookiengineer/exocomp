package shell

import "crypto/rand"
import "encoding/hex"

func (tool *Shell) newSessionID() string {

	for attempt := 0; attempt < 32; attempt++ {

		buffer := make([]byte, 4)
		_, err := rand.Read(buffer)

		if err == nil {

			session_id := "sh_" + hex.EncodeToString(buffer)

			tool.mutex.Lock()
			_, exists := tool.sessions[session_id]
			tool.mutex.Unlock()

			if exists == false {
				return session_id
			}

		} else {
			return ""
		}

	}

	return ""

}
