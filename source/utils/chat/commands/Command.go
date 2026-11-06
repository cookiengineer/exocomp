package commands

import "exocomp/schemas"
import "sort"
import "strings"

type Command struct {
	Name      string
	Method    string
	Arguments map[string]any
}

type Catalog struct {
	tools  []schemas.Tool
	byName map[string]schemas.Tool
}

func NewCatalog(tools []schemas.Tool) *Catalog {

	catalog := &Catalog{
		tools:  make([]schemas.Tool, 0),
		byName: make(map[string]schemas.Tool),
	}

	for _, tool := range tools {

		name := strings.TrimSpace(tool.Function.Name)

		if name != "" {

			catalog.tools = append(catalog.tools, tool)
			catalog.byName[strings.ToLower(name)] = tool

		}

	}

	sort.Slice(catalog.tools, func(a int, b int) bool {
		return catalog.tools[a].Function.Name < catalog.tools[b].Function.Name
	})

	return catalog

}

func (catalog *Catalog) Tools() []schemas.Tool {
	return catalog.tools
}

func (catalog *Catalog) Tool(name string) (schemas.Tool, bool) {

	tool, ok := catalog.byName[strings.ToLower(strings.TrimSpace(name))]

	return tool, ok

}

func (catalog *Catalog) Names() []string {

	result := make([]string, 0)

	for _, tool := range catalog.tools {
		result = append(result, tool.Function.Name)
	}

	return result

}
