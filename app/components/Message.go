//go:build wasm

// Package components contains the custom Gooey components of the Exocomp App.
package components

import "exocomp-app/schemas"
import utils_html "exocomp-app/utils/html"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"

// Message is a leaf component that renders a single chat message as `<article data-role="...">`.
type Message struct {
	Role      string                `json:"role"`
	Content   string                `json:"content"`
	Debug     bool                  `json:"debug"`
	Key       string                `json:"key"`
	Component *components.Component `json:"component"`
}

func NewMessage(message *schemas.Message, debug bool) *Message {

	element := dom.GetDocument().CreateElement("article")
	component := components.NewComponent(element)

	self := Message{
		Component: &component,
	}

	self.SetMessage(message, debug)

	return &self

}

func ToMessage(element *dom.Element) *Message {

	var self Message

	component := components.NewComponent(element)

	self.Component = &component
	self.Role = element.GetAttribute("data-role")
	self.Key = element.GetAttribute("data-key")

	return &self

}

func (message *Message) SetMessage(value *schemas.Message, debug bool) {

	if value == nil {
		return
	}

	message.Role = value.Role
	message.Debug = debug
	message.Content = utils_html.RenderMessage(value, debug)

}

func (message *Message) SetKey(key string) {
	message.Key = key
}

func (message *Message) Disable() bool {
	return false
}

func (message *Message) Enable() bool {
	return false
}

func (message *Message) SetScheduler(scheduler interfaces.Scheduler) {

	if message.Component != nil {
		message.Component.SetScheduler(scheduler)
	}

}

func (message *Message) Mount() bool {
	return message.Component.Element != nil
}

func (message *Message) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) == 1 {

		if message.Component.Element != nil {

			if utils.MatchesQuery(message.Component.Element, selectors[0]) == true {
				return message
			}

		}

	}

	return nil

}

func (message *Message) Render() *dom.Element {

	if message.Component.Element != nil {

		message.Component.Element.SetAttribute("data-role", message.Role)

		if message.Key != "" {
			message.Component.Element.SetAttribute("data-key", message.Key)
		}

		message.Component.Element.SetInnerHTML(message.Content)

	}

	return message.Component.Element

}

func (message *Message) String() string {

	html := "<article"

	if message.Role != "" {
		html += " data-role=\"" + message.Role + "\""
	}

	if message.Key != "" {
		html += " data-key=\"" + message.Key + "\""
	}

	html += ">"
	html += message.Content
	html += "</article>"

	return html

}

func (message *Message) Unmount() bool {
	return true
}
