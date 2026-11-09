package commands

import "exocomp/schemas"

func isRequired(tool schemas.Tool, key string) bool {

	for _, required := range tool.Function.Parameters.Required {

		if required == key {
			return true
		}

	}

	return false

}
