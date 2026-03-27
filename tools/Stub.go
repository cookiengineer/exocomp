package tools

import "fmt"
import "strings"

type Stub struct {
	Name  string
	Error error
}

func NewStub(name string, err error) *Stub {

	return &Stub{
		Name:  name,
		Error: err,
	}

}

func (tool *Stub) Call() (string, error) {
	return "", fmt.Errorf("#!tool:%s.???: %s", tool.Name, tool.Error)
}

func (tool *Stub) Parse(text string) (Tool, [2]int, error) {
	lines := strings.Split(text, "\n")
	return Tool(tool), [2]int{0, len(lines)}, nil
}
