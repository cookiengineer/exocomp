package schemas

import "encoding/json"

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
