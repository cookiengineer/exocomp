package programs

import "fmt"
import "slices"
import "strings"

func (tool *Programs) List() (string, error) {

	names := make([]string, 0)

	for _, name := range tool.AllowedPrograms {
		names = append(names, name)
	}

	slices.Sort(names)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("programs.List:"))

	for _, name := range names {
		result = append(result, fmt.Sprintf("Name: %s", name))
	}

	return strings.Join(result, "\n"), nil

}
