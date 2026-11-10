//go:build wasm

package components

import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"

type ChatFooter struct {
	Component *components.Component  `json:"component"`
	Content   []interfaces.Component `json:"content"`
}

func NewChatFooter() ChatFooter {

	element := dom.GetDocument().CreateElement("chat-footer")
	component := components.NewComponent(element)

	self := ChatFooter{
		Component: &component,
		Content:   make([]interfaces.Component, 0),
	}

	return self

}

func ToChatFooter(element *dom.Element) *ChatFooter {

	var self ChatFooter

	component := components.NewComponent(element)

	self.Component = &component
	self.Content = make([]interfaces.Component, 0)

	return &self

}

func (footer *ChatFooter) Disable() bool {

	var result bool

	for _, content := range footer.Content {
		if content.Disable() == true {
			result = true
		}
	}

	return result

}

func (footer *ChatFooter) Enable() bool {

	var result bool

	for _, content := range footer.Content {
		if content.Enable() == true {
			result = true
		}
	}

	return result

}

func (footer *ChatFooter) Invalidate() {
	footer.Component.InvalidateAs(footer)
}

func (footer *ChatFooter) SetScheduler(scheduler interfaces.Scheduler) {

	if footer.Component != nil {
		footer.Component.SetScheduler(scheduler)
	}

	for _, content := range footer.Content {
		content.SetScheduler(scheduler)
	}

}

func (footer *ChatFooter) Mount() bool {

	if footer.Component == nil || footer.Component.Element == nil {
		return false
	}

	elements := footer.Component.Element.Children()
	content := make([]interfaces.Component, 0)

	for _, element := range elements {

		switch element.TagName {

		case "CHAT-PROMPT":
			content = append(content, ToChatPrompt(element))

		default:
			tmp := components.NewComponent(element)
			content = append(content, &tmp)

		}

	}

	footer.Content = content

	for _, child := range footer.Content {
		child.Mount()
	}

	return true

}

func (footer *ChatFooter) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if footer.Component.Element != nil {

			if utils.MatchesQuery(footer.Component.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range footer.Content {

					tmp_component := content.Query(tmp_query)

					if tmp_component != nil {
						return tmp_component
					}

				}

			}

		}

	} else if len(selectors) == 1 {

		if footer.Component.Element != nil {

			if utils.MatchesQuery(footer.Component.Element, selectors[0]) == true {
				return footer
			}

		}

	}

	return nil

}

func (footer *ChatFooter) Render() *dom.Element {

	if footer.Component.Element != nil {
		components.ReconcileComponents(footer.Component.Element, footer.Content)
	}

	return footer.Component.Element

}

func (footer *ChatFooter) String() string {

	html := "<chat-footer>"

	for _, content := range footer.Content {
		html += content.String()
	}

	html += "</chat-footer>"

	return html

}

func (footer *ChatFooter) Unmount() bool {

	for _, content := range footer.Content {
		content.Unmount()
	}

	return true

}
