package components

type Builtin struct {
	Name        string
	Description string
}

var Builtins = []Builtin{
	{Name: "/help", Description: "show available commands"},
	{Name: "/clear", Description: "clear the transcript view"},
	{Name: "/copy", Description: "copy the composer through OSC52"},
	{Name: "/sidebar", Description: "toggle the changelog/bug sidebar"},
	{Name: "/quit", Description: "exit exocomp"},
}
