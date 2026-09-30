//go:build wasm

package components

import app_fmt "exocomp-app/utils/fmt"
import app_schemas "exocomp-app/schemas"
import utils_html "exocomp-app/utils/html"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import "sort"
import "strings"

type ToolCall struct {
	Component   *components.Component `json:"component"`
	Tools       []app_schemas.Tool    `json:"-"`
	Prompt      string                `json:"prompt"`
	Suggestions string                `json:"suggestions"`
	list        *dom.Element          `json:"-"`
}

type ToolCallSuggestion struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

func NewToolCall() ToolCall {

	element := dom.GetDocument().CreateElement("tool-call")
	component := components.NewComponent(element)

	self := ToolCall{
		Component:   &component,
		Tools:       make([]app_schemas.Tool, 0),
		Suggestions: "",
	}

	return self

}

func ToToolCall(element *dom.Element) *ToolCall {

	var self ToolCall

	component := components.NewComponent(element)

	self.Component = &component
	self.Tools = make([]app_schemas.Tool, 0)
	self.Suggestions = ""

	return &self

}

func (popover *ToolCall) SetTools(tools []app_schemas.Tool) {

	if tools != nil {
		popover.Tools = make([]app_schemas.Tool, len(tools))
		copy(popover.Tools, tools)
	}

}

func (popover *ToolCall) SetPrompt(prompt string, suggest bool) {

	popover.Prompt = prompt

	if strings.HasPrefix(prompt, "/") == false {

		popover.Suggestions = ""

		popover.Render()

		return

	}

	tool_command := strings.Split(prompt, " ")[0]

	tool_parameters := make([]string, 0)

	if strings.Contains(prompt, " ") == true {
		tool_parameters = strings.Split(prompt, " ")[1:]
	}

	type labeled struct {
		label       string
		description string
	}

	tool_suggestions := make([]labeled, 0)

	for _, schema := range popover.Tools {

		label := "/" + schema.Function.Name

		if strings.HasPrefix(label, tool_command) == true {
			tool_suggestions = append(tool_suggestions, labeled{label: label, description: schema.Function.Description})
		}

	}

	if len(tool_suggestions) > 1 {

		items := make([]string, 0)

		for _, suggestion := range tool_suggestions {
			items = append(items, utils_html.RenderSuggestion(map[string]any{
				"label":       suggestion.label,
				"description": suggestion.description,
			}))
		}

		sort.Strings(items)

		popover.Suggestions = strings.Join(items, "")

		popover.Render()

		return

	}

	if len(tool_suggestions) == 1 {

		items := []string{utils_html.RenderSuggestion(map[string]any{
			"label":       tool_suggestions[0].label,
			"description": tool_suggestions[0].description,
		})}

		name := strings.TrimPrefix(tool_command, "/")

		var schema *app_schemas.Tool = nil

		for t := 0; t < len(popover.Tools); t++ {

			if popover.Tools[t].Function.Name == name {
				schema = &popover.Tools[t]
				break
			}

		}

		var next *ToolCallSuggestion = nil

		if schema != nil {

			required := make(map[string]bool, len(schema.Function.Parameters.Required))
			keys := make([]string, 0)

			for _, key := range schema.Function.Parameters.Required {
				required[key] = true
				keys = append(keys, key)
			}

			optional := make([]string, 0)

			for key := range schema.Function.Parameters.Properties {

				if required[key] == false {
					optional = append(optional, key)
				}

			}

			sort.Strings(optional)

			keys = append(keys, optional...)

			for _, key := range keys {

				property, ok := schema.Function.Parameters.Properties[key]

				if ok == false {
					continue
				}

				suggestion := buildSuggestion(key, property)

				if suggestion == nil {
					continue
				}

				if next == nil && required[key] == true {

					found := false

					for _, parameter := range tool_parameters {

						if parameter != "" && strings.HasPrefix(parameter, suggestion.Key+"=") == true {
							found = true
							break
						}

					}

					if found == false {
						next = suggestion
					}

				}

				items = append(items, utils_html.RenderSuggestion(suggestion.Map()))

			}

		}

		popover.Suggestions = strings.Join(items, "")

		if next != nil && suggest == true && popover.Component != nil {
			popover.Component.FireEventListeners("suggest", next.Map())
		}

		popover.Render()

		return

	}

	popover.Suggestions = ""

	popover.Render()

}

func (popover *ToolCall) Show() {

	if popover.Component.Element != nil && popover.Component.Element.Value != nil {
		popover.Component.Element.Value.Call("showPopover")
	}

}

func (popover *ToolCall) Hide() {

	if popover.Component.Element != nil && popover.Component.Element.Value != nil {
		popover.Component.Element.Value.Call("hidePopover")
	}

}

func (popover *ToolCall) IsVisible() bool {

	if popover.Component.Element == nil || popover.Component.Element.Value == nil {
		return false
	}

	return popover.Component.Element.Value.Call("matches", ":popover-open").Bool()

}

func (popover *ToolCall) Disable() bool {
	return false
}

func (popover *ToolCall) Enable() bool {
	return false
}

func (popover *ToolCall) SetScheduler(scheduler interfaces.Scheduler) {

	if popover.Component != nil {
		popover.Component.SetScheduler(scheduler)
	}

}

func (popover *ToolCall) Mount() bool {

	if popover.Component.Element != nil {

		popover.Component.InitEvent("suggest")

		popover.Component.Element.SetAttribute("popover", "")

		popover.list = popover.Component.Element.QuerySelector("ul[data-name=\"suggestions\"]")

		return true

	}

	return false

}

func (popover *ToolCall) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if popover.Component.Element != nil {

			if utils.MatchesQuery(popover.Component.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range popover.Component.Content {

					tmp_component := content.Query(tmp_query)

					if tmp_component != nil {
						return tmp_component
					}

				}

			}

		}

	} else if len(selectors) == 1 {

		if popover.Component.Element != nil {

			if utils.MatchesQuery(popover.Component.Element, selectors[0]) == true {
				return popover
			}

		}

	}

	return nil

}

func (popover *ToolCall) Render() *dom.Element {

	if popover.Component.Element != nil {

		popover.Component.Element.SetAttribute("popover", "")

		if popover.list != nil {
			popover.list.SetInnerHTML(popover.Suggestions)
		}

		if strings.HasPrefix(popover.Prompt, "/") == true {

			if popover.IsVisible() == false {
				popover.Show()
			}

		} else {

			if popover.IsVisible() == true {
				popover.Hide()
			}

		}

	}

	return popover.Component.Element

}

func (popover *ToolCall) String() string {

	html := "<tool-call>"

	if popover.Component.Element != nil {
		html += popover.Component.Element.InnerHTML
	}

	html += "</tool-call>"

	return html

}

func (popover *ToolCall) Unmount() bool {

	if popover.Component.Element != nil {
		popover.Component.RemoveEventListener("suggest", nil)
	}

	return true

}

func (suggestion *ToolCallSuggestion) Map() map[string]any {

	return map[string]any{
		"key":         suggestion.Key,
		"type":        suggestion.Type,
		"value":       suggestion.Value,
		"label":       suggestion.Label,
		"description": suggestion.Description,
	}

}

func buildSuggestion(key string, property app_schemas.ToolFunctionParameterProperty) *ToolCallSuggestion {

	description := property.Description

	if property.Type == "array" {

		if property.Items == nil {
			return nil
		}

		switch property.Items.Type {

		case "boolean":
			return &ToolCallSuggestion{
				Key:         key,
				Type:        "array-of-booleans",
				Value:       "[true,...]",
				Label:       key + "=[true,...]",
				Description: description,
			}

		case "number":
			return &ToolCallSuggestion{
				Key:         key,
				Type:        "array-of-numbers",
				Value:       "[13,37,...]",
				Label:       key + "=[13,37,...]",
				Description: description,
			}

		case "string":
			return &ToolCallSuggestion{
				Key:         key,
				Type:        "array-of-strings",
				Value:       "[\"foo\",...]",
				Label:       key + "=[\"foo\",...]",
				Description: description,
			}

		}

		return nil

	} else if property.Type == "boolean" {

		return &ToolCallSuggestion{
			Key:         key,
			Type:        "boolean",
			Value:       "true",
			Label:       key + "=true",
			Description: description,
		}

	} else if property.Type == "number" {

		return &ToolCallSuggestion{
			Key:         key,
			Type:        "number",
			Value:       "1337",
			Label:       key + "=1337",
			Description: description,
		}

	} else if property.Type == "string" {

		if len(property.Enum) > 0 {

			value := "\"" + property.Enum[0] + "\""

			return &ToolCallSuggestion{
				Key:         key,
				Type:        "string",
				Value:       value,
				Label:       key + "=" + value,
				Description: description + " (" + app_fmt.JoinQuoted(property.Enum) + ")",
			}

		}

		return &ToolCallSuggestion{
			Key:         key,
			Type:        "string",
			Value:       "\"...\"",
			Label:       key + "=\"...\"",
			Description: description,
		}

	}

	return nil

}
