package shell

import "fmt"

func (tool *Shell) Resize(session_id string, columns int, rows int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if columns >= shellMinColumns && columns <= shellMaxColumns {

			if rows >= shellMinRows && rows <= shellMaxRows {

				err0 := session.resize(columns, rows)

				if err0 == nil {
					return fmt.Sprintf("shell.Resize: Session \"%s\" resized to %dx%d.", session_id, columns, rows), nil
				} else {
					return "", fmt.Errorf("shell.Resize: Cannot resize session \"%s\".", session_id)
				}

			} else {
				return "", fmt.Errorf("shell.Resize: Invalid parameter \"rows\" must be between %d and %d.", shellMinRows, shellMaxRows)
			}

		} else {
			return "", fmt.Errorf("shell.Resize: Invalid parameter \"columns\" must be between %d and %d.", shellMinColumns, shellMaxColumns)
		}

	} else {
		return "", fmt.Errorf("shell.Resize: Session \"%s\" not found.", session_id)
	}

}
