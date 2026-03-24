package tools

type Tool interface {
	Call()        (string, error)
	Parse(string) (Tool, [2]int, error)
}
