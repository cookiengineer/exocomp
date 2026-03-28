package schemas

import "encoding/json"
import "fmt"
import "strings"

type ToolCall struct {
	ID       string           `json:"id,omitempty"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name         string          `json:"name"`
	ArgumentsRaw json.RawMessage `json:"arguments"` // some models return JSON string
	// Arguments map[string]interface{} `json:"arguments"`
}

func (function *ToolCallFunction) Tool() (string, error) {

	if strings.Contains(function.Name, ".") {

		tmp := strings.Split(function.Name, ".")

		if len(tmp) == 2 && len(tmp[0]) > 0 {

			tool := strings.TrimSpace(strings.ToLower(tmp[0]))

			if tool != "" {
				return tool, nil
			} else {
				return "", fmt.Errorf("Invalid Tool Call")
			}

		} else {
			return "", fmt.Errorf("Invalid Tool Call")
		}

	} else {
		return "", fmt.Errorf("Invalid Tool Call")
	}

}

func (function *ToolCallFunction) Method() (string, error) {

	if strings.Contains(function.Name, ".") {

		tmp := strings.Split(function.Name, ".")

		if len(tmp) == 2 && len(tmp[0]) > 0 && len(tmp[1]) >= 2 {

			method := strings.TrimSpace(strings.ToUpper(tmp[1][0:1]) + strings.ToLower(tmp[1][1:]))

			if method != "" {
				return method, nil
			} else {
				return "", fmt.Errorf("Invalid Tool Call")
			}

		} else {
			return "", fmt.Errorf("Invalid Tool Call")
		}

	} else {
		return "", fmt.Errorf("Invalid Tool Call")
	}

}

func (function *ToolCallFunction) Arguments() (map[string]interface{}, error) {

	result := make(map[string]interface{})
	err0   := json.Unmarshal(function.ArgumentsRaw, &result)

	if err0 == nil {

		// Arguments was an Object
		return result, nil

	} else {

		// Arguments was a JSON string
		encoded := ""
		err1    := json.Unmarshal(function.ArgumentsRaw, &encoded)

		if err1 == nil {

			err2 := json.Unmarshal([]byte(encoded), &result)

			if err2 == nil {
				return result, nil
			} else {
				return nil, err2
			}

		} else {
			return nil, err1
		}

	}

}
