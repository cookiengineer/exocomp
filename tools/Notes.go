package tools

import "fmt"

type Notes struct {
	Method    string
	Arguments []string
	Sandbox   string
}

func NewNotes(sandbox string, tools []string, programs []string) *Notes {

	return &Notes{
		Method:    "",
		Arguments: make([]string, 0),
		Sandbox:   sandbox,
	}

}

func (tool *Notes) Call() (string, error) {

	if tool.Method == "Help" {
	} else if tool.Method == "Create" {
	} else if tool.Method == "List" {
	} else {
		return "", fmt.Errorf("#!tool:notes.%s: Unknown method", tool.Method)
	}

}

// TODO: Rewrite tasks API to use NOTES.md file instead

