package tools

import "fmt"

type Requirements struct {
	Sandbox string
}

func NewRequirements(agent string, sandbox string) *Requirements {

	requirements := &Requirements{
		Sandbox: sandbox,
	}

	return requirements

}

func (tool *Requirements) Call(method string, arguments map[string]interface{}) (string, error) {
	return "", fmt.Errorf("agents.%s: Not implemented.", method)
}
