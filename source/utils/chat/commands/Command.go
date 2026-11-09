package commands

type Command struct {
	Name      string
	Method    string
	Arguments map[string]any
}
