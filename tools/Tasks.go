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
	} else if tool.Method == "List" {
	} else if tool.Method == "Start" {
	} else if tool.Method == "Pause" {
	} else if tool.Method == "Stop" {
	} else if tool.Method == "Reopen" {
	} else {
		return "", fmt.Errorf("#!tool:tasks.%s: Unknown method", tool.Method)
	}

}

// TODO: Rewrite tasks API to use TODO.md file instead
