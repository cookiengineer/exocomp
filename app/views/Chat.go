//go:build wasm

package views

import app_components "exocomp-app/components"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/layout"
import "github.com/cookiengineer/gooey/components/types"
import "github.com/cookiengineer/gooey/components/utils"
import "strings"

type Chat struct {
	Element   *dom.Element           `json:"element"`
	Layout    types.Layout           `json:"layout"`
	Content   []interfaces.Component `json:"content"`
	name      string
	label     string
	path      string
	scheduler interfaces.Scheduler
}

func ToChat(element *dom.Element) *Chat {

	var view Chat

	view.Element = element
	view.Layout = types.LayoutFlow
	view.Content = make([]interfaces.Component, 0)

	view.name = strings.ToLower(element.GetAttribute("data-name"))
	view.label = element.GetAttribute("data-label")
	view.path = strings.ToLower(element.GetAttribute("data-path"))

	return &view

}

func (view *Chat) Disable() bool {
	return false
}

func (view *Chat) Enable() bool {
	return false
}

func (view *Chat) Enter() bool {

	if view.Element != nil {
		view.Element.SetAttribute("data-state", "active")
	}

	return true

}

func (view *Chat) Leave() bool {

	if view.Element != nil {
		view.Element.RemoveAttribute("data-state")
	}

	return true

}

func (view *Chat) Label() string {
	return view.label
}

func (view *Chat) Name() string {
	return view.name
}

func (view *Chat) Path() string {
	return view.path
}

func (view *Chat) SetScheduler(scheduler interfaces.Scheduler) {

	view.scheduler = scheduler

	for _, component := range view.Content {

		if component != nil {
			component.SetScheduler(scheduler)
		}

	}

}

func (view *Chat) Mount() bool {

	if view.Element == nil {
		return false
	}

	tmp_name := view.Element.GetAttribute("data-name")
	if tmp_name != "" {
		view.name = strings.ToLower(tmp_name)
	}

	tmp_label := view.Element.GetAttribute("data-label")
	if tmp_label != "" {
		view.label = tmp_label
	}

	tmp_layout := view.Element.GetAttribute("data-layout")
	if tmp_layout != "" {
		view.Layout = types.Layout(tmp_layout)
	}

	tmp_path := view.Element.GetAttribute("data-path")
	if tmp_path != "" {
		view.path = strings.ToLower(tmp_path)
	}

	elements := view.Element.Children()
	content := make([]interfaces.Component, 0)

	for _, element := range elements {

		switch element.TagName {

		case "ASIDE", "AGENT-LIST":
			content = append(content, app_components.ToAgentList(element))

		case "MESSAGE-LIST":
			content = append(content, app_components.ToMessageList(element))

		case "CHAT-FOOTER":
			content = append(content, app_components.ToChatFooter(element))

		case "DIALOG":
			content = append(content, layout.ToDialog(element))

		default:
			component := components.NewComponent(element)
			content = append(content, &component)

		}

	}

	view.Content = content

	for _, component := range view.Content {
		component.Mount()
	}

	return true

}

func (view *Chat) Unmount() bool {

	for _, component := range view.Content {
		component.Unmount()
	}

	return true

}

// GetAgentList returns the mapped agents sidebar component.
func (view *Chat) GetAgentList() *app_components.AgentList {

	for _, component := range view.Content {

		if component == nil {
			continue
		}

		if agent_list, ok := component.(*app_components.AgentList); ok == true {
			return agent_list
		}

	}

	return nil

}

// GetMessageList returns the mapped message list component.
func (view *Chat) GetMessageList() *app_components.MessageList {

	for _, component := range view.Content {

		if component == nil {
			continue
		}

		if message_list, ok := component.(*app_components.MessageList); ok == true {
			return message_list
		}

	}

	return nil

}

func (view *Chat) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if view.Element != nil {

			if utils.MatchesQuery(view.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range view.Content {

					tmp_component := content.Query(tmp_query)

					if tmp_component != nil {
						return tmp_component
					}

				}

			}

		}

	} else if len(selectors) == 1 {

		if view.Element != nil {

			if utils.MatchesQuery(view.Element, selectors[0]) == true {
				return view
			}

		}

	}

	return nil

}

func (view *Chat) QuerySelector(query string) *dom.Element {

	if view.Element != nil {
		return view.Element.QuerySelector(query)
	}

	return nil

}

func (view *Chat) QuerySelectorAll(query string) []*dom.Element {

	result := make([]*dom.Element, 0)

	if view.Element != nil {
		result = view.Element.QuerySelectorAll(query)
	}

	return result

}

func (view *Chat) Render() *dom.Element {

	if view.Element != nil {

		if view.name != "" {
			view.Element.SetAttribute("data-name", view.name)
		}

		if view.label != "" {
			view.Element.SetAttribute("data-label", view.label)
		}

		if view.path != "" {
			view.Element.SetAttribute("data-path", view.path)
		}

		if view.Layout != types.LayoutFlow {
			view.Element.SetAttribute("data-layout", view.Layout.String())
		}

		components.ReconcileComponents(view.Element, view.Content)

		return view.Element

	}

	return nil

}

func (view *Chat) String() string {

	html := "<section"

	if view.name != "" {
		html += " data-name=\"" + view.name + "\""
	}

	if view.label != "" {
		html += " data-label=\"" + view.label + "\""
	}

	if view.path != "" {
		html += " data-path=\"" + view.path + "\""
	}

	html += ">"

	for _, content := range view.Content {
		html += content.String()
	}

	html += "</section>"

	return html

}
