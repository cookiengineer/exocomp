package schemas

import "encoding/json"
import "fmt"
import "strings"

func toUpperCase(raw string) string {

	check := strings.ToLower(raw)

	// NOTE: Dirty Hack for stupid LLMs
	if strings.HasPrefix(check, "define") {

		tmp2 := strings.TrimSpace(raw[6:])

		return "Define" + strings.ToUpper(tmp2[0:1]) + strings.ToLower(tmp2[1:])

	} else {
		return strings.ToUpper(raw[0:1]) + strings.ToLower(raw[1:])
	}

}

type ToolCall struct {
	ID       string           `json:"id,omitempty" yaml:"id,omitempty"`
	Type     string           `json:"type" yaml:"type"`
	Function ToolCallFunction `json:"function" yaml:"function"`
}

type ToolCallFunction struct {
	Name         string          `json:"name" yaml:"name"`
	ArgumentsRaw json.RawMessage `json:"arguments" yaml:"arguments"` // some models return JSON string
	// Arguments map[string]interface{} `json:"arguments"`
}

func (toolcall *ToolCall) ToolID() (string, error) {
	return strings.TrimSpace(toolcall.ID), nil
}

func (toolcall *ToolCall) ToolName() (string, error) {

	if strings.Contains(toolcall.Function.Name, ".") {

		tmp := strings.Split(toolcall.Function.Name, ".")

		if len(tmp) == 2 && len(tmp[0]) > 0 && len(tmp[1]) >= 2 {

			tool   := strings.TrimSpace(strings.ToLower(tmp[0]))
			method := strings.TrimSpace(toUpperCase(tmp[1]))

			if tool != "" && method != "" {
				return tool + "." + method, nil
			} else {
				return "", fmt.Errorf("Invalid Tool Call: Name \"%s\" is invalid.", toolcall.Function.Name)
			}

		} else {
			return "", fmt.Errorf("Invalid Tool Call: Name \"%s\" is invalid.", toolcall.Function.Name)
		}

	} else {
		return "", fmt.Errorf("Invalid Tool Call: Name \"%s\" is invalid.", toolcall.Function.Name)
	}

}

func (toolcall *ToolCall) ToolMethod() (string, error) {

	if strings.Contains(toolcall.Function.Name, ".") {

		tmp := strings.Split(toolcall.Function.Name, ".")

		if len(tmp) == 2 && len(tmp[0]) > 0 && len(tmp[1]) >= 2 {

			method := strings.TrimSpace(toUpperCase(tmp[1]))

			if method != "" {
				return method, nil
			} else {
				return "", fmt.Errorf("Invalid Tool Call: Method from \"%s\" is invalid.", toolcall.Function.Name)
			}

		} else {
			return "", fmt.Errorf("Invalid Tool Call: Method from \"%s\" is invalid.", toolcall.Function.Name)
		}

	} else {
		return "", fmt.Errorf("Invalid Tool Call: Method from \"%s\" is invalid.", toolcall.Function.Name)
	}

}

func (toolcall *ToolCall) ToolArguments() (map[string]any, error) {

	result := make(map[string]any)
	err0   := json.Unmarshal(toolcall.Function.ArgumentsRaw, &result)

	if err0 == nil {

		// Arguments was an Object
		return result, nil

	} else {

		// Arguments was a JSON string
		encoded := ""
		err1    := json.Unmarshal(toolcall.Function.ArgumentsRaw, &encoded)

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
