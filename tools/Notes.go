package tools

import "fmt"

type Notes struct {
	Method    string
	Arguments []string
	Sandbox   string
	database  map[string]string
}

func NewNotes(sandbox string, tools []string, programs []string) *Notes {

	notes := &Notes{
		Method:    "",
		Arguments: make([]string, 0),
		Sandbox:   sandbox,
		database:  make(map[string]string),
	}

	readNotes(notes)

	return notes

}

func (tool *Notes) Call() (string, error) {

	if tool.Method == "Help" {
		return tool.Help(tool.Arguments)
	} else if tool.Method == "Create" {
		return tool.Create(tool.Arguments)
	} else if tool.Method == "List" {
		return tool.List(tool.Arguments)
	} else {
		return "", fmt.Errorf("#!tool:notes.%s: Unknown method", tool.Method)
	}

}

func (tool *Notes) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!tool:notes.Create \"Note description with whitespaces\"",
		"",
		"#!tool:notes.List",
		"",
		"#!tool:notes.Search \"keyword\"",
		"",
		"#!tool:notes.Remove note-id",
	}, "\n"), nil

}

func (tool *Notes) Create(arguments []string) (string, error) {

	if len(arguments) == 1 {

		description := strings.TrimSpace(arguments[0])

		if description != "" {
		} else {
		}

		// TODO

	} else {
		return "", fmt.Errorf("#!tool:notes.Create: Only one argument allowed")
	}

}

func (tool *Notes) List(arguments []string) (string, error) {

	if len(arguments) == 0 {

		// TODO

	} else {
		return "", fmt.Errorf("#!tool:notes.List: Only zero arguments allowed")
	}

}

func (tool *Notes) Search(arguments []string) (string, error) {

	if len(arguments) == 1 {

		// TODO

	} else {
		return "", fmt.Errorf("#!tool:notes.Search: Only one argument allowed")
	}

}

func (tool *Notes) Remove(arguments []string) (string, error) {

	if len(arguments) == 1 {

		// TODO

	} else {
		return "", fmt.Errorf("#!tool:notes.Remove: Only one argument allowed")
	}

}

