package types

import "encoding/json"
import "strconv"

func parseCommandArgument(buffer []byte) any {

	var value any

	err := json.Unmarshal(buffer, &value)

	if err == nil {

		return value

	} else {

		if len(buffer) >= 2 && (buffer[0] == '-' || buffer[0] == '+') {

			num, err1 := strconv.ParseInt(string(buffer), 10, 64)

			if err1 == nil {
				return num
			}

		}

		if len(buffer) >= 2 && buffer[0] == '"' && buffer[len(buffer)-1] == '"' {
			return string(buffer[1 : len(buffer)-1])
		}

		if len(buffer) >= 2 && buffer[0] == '\'' && buffer[len(buffer)-1] == '\'' {
			return string(buffer[1 : len(buffer)-1])
		}

		if string(buffer) == "true" {
			return true
		}

		if string(buffer) == "false" {
			return false
		}

		if string(buffer) == "null" {
			return nil
		}

		num1, err1 := strconv.ParseUint(string(buffer), 10, 64)

		if err1 == nil {
			return num1
		}

		num2, err2 := strconv.ParseFloat(string(buffer), 64)

		if err2 == nil {
			return num2
		}

	}

	return nil

}


