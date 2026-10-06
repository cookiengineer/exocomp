//go:build wasm

package components

import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import std_html "html"
import "strings"

// BugFooter is the bottom bar of the Bugs view. It owns the search field and the
// collapse toggle and emits `toggle-collapse` / `change-filter` events.
type BugFooter struct {
	Component   *components.Component `json:"component"`
	Collapsed   bool                  `json:"collapsed"`
	Placeholder string                `json:"placeholder"`
	button      *dom.Element          `json:"-"`
	search      *dom.Element          `json:"-"`
	initialized bool                  `json:"-"`

	on_click *dom.EventListener `json:"-"`
	on_input *dom.EventListener `json:"-"`
}

func NewBugFooter(placeholder string) *BugFooter {

	element := dom.GetDocument().CreateElement("footer")
	component := components.NewComponent(element)

	return &BugFooter{
		Component:   &component,
		Placeholder: placeholder,
	}

}

func (footer *BugFooter) SetCollapsed(collapsed bool) {

	footer.Collapsed = collapsed

	footer.renderButton()

}

func (footer *BugFooter) Disable() bool {
	return false
}

func (footer *BugFooter) Enable() bool {
	return false
}

func (footer *BugFooter) SetScheduler(scheduler interfaces.Scheduler) {

	if footer.Component != nil {
		footer.Component.SetScheduler(scheduler)
	}

}

func (footer *BugFooter) ensure() {

	if footer.initialized == true {
		return
	}

	if footer.Component == nil || footer.Component.Element == nil {
		return
	}

	placeholder := strings.TrimSpace(footer.Placeholder)

	if placeholder == "" {
		placeholder = "Search ..."
	}

	html := "<div></div>"
	html += "<div><input type=\"search\" data-name=\"search\" placeholder=\"" + std_html.EscapeString(placeholder) + "\"/></div>"
	html += "<div><button type=\"button\" data-action=\"toggle-collapse\" data-state=\"nested\">Collapse Packages</button></div>"

	footer.Component.Element.SetInnerHTML(html)

	footer.button = footer.Component.Element.QuerySelector("button[data-action=\"toggle-collapse\"]")
	footer.search = footer.Component.Element.QuerySelector("input[data-name=\"search\"]")

	footer.on_click = dom.ToEventListener(func(event *dom.Event) {

		if event.Target == nil {
			return
		}

		if event.Target.GetAttribute("data-action") == "toggle-collapse" {
			footer.Component.FireEventListeners("toggle-collapse", map[string]any{})
		}

	})

	footer.Component.Element.AddEventListener(dom.EventType("click"), footer.on_click)

	footer.on_input = dom.ToEventListener(func(event *dom.Event) {

		if footer.search == nil || footer.search.Value == nil {
			return
		}

		value := footer.search.Value.Get("value")

		if value.IsNull() || value.IsUndefined() {
			return
		}

		footer.Component.FireEventListeners("change-filter", map[string]any{"value": value.String()})

	})

	footer.Component.Element.AddEventListener(dom.EventType("input"), footer.on_input)

	footer.initialized = true

}

func (footer *BugFooter) renderButton() {

	if footer.button == nil {
		return
	}

	if footer.Collapsed == true {
		footer.button.SetAttribute("data-state", "collapsed")
		footer.button.SetTextContent("Expand Packages")
	} else {
		footer.button.SetAttribute("data-state", "nested")
		footer.button.SetTextContent("Collapse Packages")
	}

}

func (footer *BugFooter) Mount() bool {

	if footer.Component == nil || footer.Component.Element == nil {
		return false
	}

	footer.Component.InitEvent("toggle-collapse")
	footer.Component.InitEvent("change-filter")

	footer.ensure()
	footer.renderButton()

	return true

}

func (footer *BugFooter) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) == 1 {

		if footer.Component != nil && footer.Component.Element != nil {

			if utils.MatchesQuery(footer.Component.Element, selectors[0]) == true {
				return footer
			}

		}

	}

	return nil

}

func (footer *BugFooter) Render() *dom.Element {

	if footer.Component == nil || footer.Component.Element == nil {
		return nil
	}

	footer.ensure()
	footer.renderButton()

	return footer.Component.Element

}

func (footer *BugFooter) String() string {

	html := "<footer>"

	if footer.Component != nil && footer.Component.Element != nil {
		html += footer.Component.Element.InnerHTML
	}

	html += "</footer>"

	return html

}

func (footer *BugFooter) Unmount() bool {

	if footer.Component != nil && footer.Component.Element != nil {

		if footer.on_click != nil {
			footer.Component.Element.RemoveEventListener(dom.EventType("click"), footer.on_click)
			footer.on_click = nil
		}

		if footer.on_input != nil {
			footer.Component.Element.RemoveEventListener(dom.EventType("input"), footer.on_input)
			footer.on_input = nil
		}

	}

	return true

}
