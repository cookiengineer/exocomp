package commands

import "exocomp/schemas"
import "sort"
import "strings"

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

func (catalog *Catalog) Complete(input string, limit int) []Suggestion {

	if strings.HasPrefix(input, "/") == false {
		return []Suggestion{}
	}

	content := strings.TrimPrefix(input, "/")

	if strings.ContainsAny(content, " \t\n") == true {

		parts := strings.Fields(content)

		if len(parts) == 0 {
			return []Suggestion{}
		}

		tool, ok := catalog.Tool(parts[0])

		if ok == false {
			return []Suggestion{}
		}

		raw, _ := ParseRawArguments(strings.Join(parts[1:], " "))
		result := make([]Suggestion, 0)
		keys := make([]string, 0)

		for key := range tool.Function.Parameters.Properties {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {

			_, has_value := raw[key]

			if has_value == false {

				property := tool.Function.Parameters.Properties[key]

				result = append(result, Suggestion{
					Kind:        "parameter",
					Value:       key,
					Description: property.Description,
					Type:        property.Type,
					Required:    isRequired(tool, key),
				})

			}

		}

		if limit > 0 && len(result) > limit {
			result = result[0:limit]
		}

		return result

	} else {

		prefix := strings.ToLower(content)
		result := make([]Suggestion, 0)

		for _, tool := range catalog.tools {

			name := tool.Function.Name

			if prefix == "" || strings.Contains(strings.ToLower(name), prefix) {

				result = append(result, Suggestion{
					Kind:        "tool",
					Value:       name,
					Description: tool.Function.Description,
				})

			}

		}

		sort.Slice(result, func(a int, b int) bool {

			a_prefix := strings.HasPrefix(strings.ToLower(result[a].Value), prefix)
			b_prefix := strings.HasPrefix(strings.ToLower(result[b].Value), prefix)

			if a_prefix != b_prefix {
				return a_prefix
			} else {
				return result[a].Value < result[b].Value
			}

		})

		if limit > 0 && len(result) > limit {
			result = result[0:limit]
		}

		return result

	}

}
