//go:build wasm

package components

import app_types "exocomp-app/types"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import "fmt"
import std_html "html"
import "net/url"
import "sort"
import "strings"

// AgentList is a layout component that renders the agents navigation as an
// `<aside>` sidebar so that the Gooey theme's `aside` styles are inherited. It
// fires a `change-agent` event when an agent is clicked, mirroring the old
// frontend `aside > nav[aria-label="agents"]` sidebar.
type AgentList struct {
	Component *components.Component       `json:"component"`
	Active    string                      `json:"active"`
	Agents    map[string]*app_types.Agent `json:"-"`
}

func NewAgentList() AgentList {

	element := dom.GetDocument().CreateElement("aside")
	component := components.NewComponent(element)

	self := AgentList{
		Component: &component,
		Agents:    make(map[string]*app_types.Agent),
	}

	return self

}

func ToAgentList(element *dom.Element) *AgentList {

	var self AgentList

	component := components.NewComponent(element)

	self.Component = &component
	self.Agents = make(map[string]*app_types.Agent)

	return &self

}

func (list *AgentList) SetAgents(active string, agents map[string]*app_types.Agent) {

	list.Active = active

	if agents != nil {
		list.Agents = agents
	}

	list.Render()

}

func (list *AgentList) Disable() bool {
	return false
}

func (list *AgentList) Enable() bool {
	return false
}

func (list *AgentList) SetScheduler(scheduler interfaces.Scheduler) {

	if list.Component != nil {
		list.Component.SetScheduler(scheduler)
	}

}

func (list *AgentList) Mount() bool {

	if list.Component.Element != nil {

		list.Component.InitEvent("change-agent")

		list.Component.Element.AddEventListener(dom.EventType("click"), dom.ToEventListener(func(event *dom.Event) {

			if event.Target == nil {
				return
			}

			name := ""

			if event.Target.TagName == "LABEL" || event.Target.TagName == "A" {
				name = strings.TrimSpace(event.Target.GetTextContent())
			}

			if name != "" {
				list.Component.FireEventListeners("change-agent", map[string]any{"name": name})
			}

			event.PreventDefault()
			event.StopPropagation()

		}))

		return true

	}

	return false

}

func (list *AgentList) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if list.Component.Element != nil {

			if utils.MatchesQuery(list.Component.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range list.Component.Content {

					tmp_component := content.Query(tmp_query)

					if tmp_component != nil {
						return tmp_component
					}

				}

			}

		}

	} else if len(selectors) == 1 {

		if list.Component.Element != nil {

			if utils.MatchesQuery(list.Component.Element, selectors[0]) == true {
				return list
			}

		}

	}

	return nil

}

func (list *AgentList) Render() *dom.Element {

	if list.Component.Element != nil {

		names := make([]string, 0)

		for name := range list.Agents {
			names = append(names, name)
		}

		sort.Strings(names)

		html := "<nav aria-label=\"agents\"><ul>"

		for _, name := range names {

			state := ""

			if name == list.Active {
				state = "active"
			}

			sandbox := ""
			title := name

			if agent := list.Agents[name]; agent != nil {

				sandbox = agent.Sandbox
				title = name + " working in " + sandbox

			}

			href := "/agent.html?name=" + url.QueryEscape(name)

			html += fmt.Sprintf("<li data-state=\"%s\" title=\"%s\"><a href=\"%s\"><label>%s</label></a></li>", state, std_html.EscapeString(title), href, std_html.EscapeString(name))

		}

		html += "</ul></nav>"

		list.Component.Element.SetInnerHTML(html)

	}

	return list.Component.Element

}

func (list *AgentList) String() string {

	html := "<aside>"

	if list.Component.Element != nil {
		html += list.Component.Element.InnerHTML
	}

	html += "</aside>"

	return html

}

func (list *AgentList) Unmount() bool {

	if list.Component.Element != nil {
		list.Component.Element.RemoveEventListener(dom.EventType("click"), nil)
	}

	return true

}
