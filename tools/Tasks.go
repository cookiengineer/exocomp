package tools

import "fmt"

type Tasks struct {
	Method    string
	Arguments []string
	Sandbox   string
}

func NewTasks(agent string, sandbox string, tools []string, programs []string) *Tasks {

	return &Tasks{
		Method:    "",
		Arguments: make([]string, 0),
		Sandbox:   sandbox,
	}

}

func (tool *Tasks) Call() (string, error) {

	if tool.Method == "Create" {
		return tool.NotImplemented(tool.Arguments)
	} else if tool.Method == "List" {
		return tool.NotImplemented(tool.Arguments)
	} else if tool.Method == "Start" {
		return tool.NotImplemented(tool.Arguments)
	} else if tool.Method == "Pause" {
		return tool.NotImplemented(tool.Arguments)
	} else if tool.Method == "Stop" {
		return tool.NotImplemented(tool.Arguments)
	} else {
		return "", fmt.Errorf("#!tool:tasks.%s: Invalid method.", tool.Method)
	}

}

func (tool *Tasks) NotImplemented(arguments []string) (string, error) {
	return "", fmt.Errorf("#!tool:tasks.NotImplemented: Method not implemented yet.")
}

// TODO: Rewrite tasks API to use TODO.md file instead
