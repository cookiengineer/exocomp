package shell

import "fmt"

func (tool *Shell) Log(session_id string, offset int, limit int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if offset >= 0 {

			if limit >= 1 && limit <= shellMaxRead {

				data, total, truncated, err0 := session.readLog(int64(offset), limit)

				if err0 == nil {

					header := fmt.Sprintf("shell.Log: Session \"%s\" log %d bytes, offset %d, read %d bytes.", session_id, total, offset, len(data))

					if truncated == true {
						header += " Log is truncated."
					}

					if len(data) > 0 {
						return header + "\n" + string(data), nil
					} else {
						return header + "\n<no data>", nil
					}

				} else {
					return "", fmt.Errorf("shell.Log: %s", err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.Log: Invalid parameter \"limit\" must be between 1 and %d.", shellMaxRead)
			}

		} else {
			return "", fmt.Errorf("shell.Log: Invalid parameter \"offset\" must be zero or positive.")
		}

	} else {
		return "", fmt.Errorf("shell.Log: Session \"%s\" not found.", session_id)
	}

}
