//go:build wasm

package components

import app_schemas "exocomp-app/schemas"
import utils_html "exocomp-app/utils/html"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import "strconv"

// MessageList is a layout component that owns the chat `Message` children. It
// mirrors the old frontend `ui/renderers/ChatRenderer.mjs` message rendering.
type MessageList struct {
	Component *components.Component  `json:"component"`
	Content   []interfaces.Component `json:"content"`
	Messages  []*app_schemas.Message `json:"-"`
	Debug     bool                   `json:"-"`
	Waiting   bool                   `json:"-"`
	thinking  *dom.Element           `json:"-"`
}

func NewMessageList() MessageList {

	element := dom.GetDocument().CreateElement("message-list")
	component := components.NewComponent(element)

	self := MessageList{
		Component: &component,
		Content:   make([]interfaces.Component, 0),
	}

	return self

}

func ToMessageList(element *dom.Element) *MessageList {

	var self MessageList

	component := components.NewComponent(element)

	self.Component = &component
	self.Content = make([]interfaces.Component, 0)

	return &self

}

func (list *MessageList) SetMessages(messages []*app_schemas.Message, debug bool) {

	list.Messages = messages
	list.Debug = debug

	content := make([]interfaces.Component, 0)

	for m := 0; m < len(messages); m++ {

		message := messages[m]

		if message == nil {
			continue
		}

		if utils_html.RenderMessage(message, debug) == "" {
			continue
		}

		child := NewMessage(message, debug)
		child.SetKey(strconv.Itoa(m))
		child.Mount()

		content = append(content, child)

	}

	list.Content = content

	list.Render()

}

func (list *MessageList) Clear() {

	list.Messages = make([]*app_schemas.Message, 0)
	list.Content = make([]interfaces.Component, 0)

	list.Render()

}

func (list *MessageList) SetWaiting(waiting bool) {

	if list.Waiting == waiting {
		return
	}

	list.Waiting = waiting
	list.Render()

}

func (list *MessageList) Disable() bool {
	return false
}

func (list *MessageList) Enable() bool {
	return false
}

func (list *MessageList) SetScheduler(scheduler interfaces.Scheduler) {

	if list.Component != nil {
		list.Component.SetScheduler(scheduler)
	}

	for _, content := range list.Content {
		content.SetScheduler(scheduler)
	}

}

func (list *MessageList) Mount() bool {

	for _, content := range list.Content {
		content.Mount()
	}

	return true

}

func (list *MessageList) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if list.Component.Element != nil {

			if utils.MatchesQuery(list.Component.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range list.Content {

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

func (list *MessageList) Render() *dom.Element {

	if list.Component.Element != nil {

		open_details := snapshotOpenDetails(list.Component.Element)

		components.ReconcileComponents(list.Component.Element, list.Content)

		restoreOpenDetails(list.Component.Element, open_details)

		if list.Waiting == true {

			if list.thinking == nil {

				list.thinking = dom.GetDocument().CreateElement("article")
				list.thinking.SetAttribute("data-role", "thinking")
				list.thinking.SetClassName("thinking")
				list.thinking.SetInnerHTML("<span class=\"spinner\"></span> Thinking")

			}

			list.Component.Element.Append(list.thinking)

		} else {
			list.thinking = nil
		}

	}

	return list.Component.Element

}

func snapshotOpenDetails(element *dom.Element) map[string]bool {

	result := make(map[string]bool)

	if element == nil {
		return result
	}

	for _, article := range element.QuerySelectorAll("article[data-key]") {

		message_key := article.GetAttribute("data-key")

		if message_key == "" {
			continue
		}

		counts := make(map[string]int)

		for _, detail := range article.QuerySelectorAll("details.message-collapse") {

			class_name := detail.GetAttribute("class")
			index := counts[class_name]
			counts[class_name] = index + 1

			if detail.HasAttribute("open") == true {
				key := message_key + ":" + class_name + ":" + strconv.Itoa(index)
				result[key] = true
			}

		}

	}

	return result

}

func restoreOpenDetails(element *dom.Element, open map[string]bool) {

	if element == nil || len(open) == 0 {
		return
	}

	for _, article := range element.QuerySelectorAll("article[data-key]") {

		message_key := article.GetAttribute("data-key")

		if message_key == "" {
			continue
		}

		counts := make(map[string]int)

		for _, detail := range article.QuerySelectorAll("details.message-collapse") {

			class_name := detail.GetAttribute("class")
			index := counts[class_name]
			counts[class_name] = index + 1

			key := message_key + ":" + class_name + ":" + strconv.Itoa(index)

			if open[key] == true {
				detail.SetAttribute("open", "")
			}

		}

	}

}

func (list *MessageList) String() string {

	html := "<message-list>"

	for _, content := range list.Content {
		html += content.String()
	}

	html += "</message-list>"

	return html

}

func (list *MessageList) Unmount() bool {

	for _, content := range list.Content {
		content.Unmount()
	}

	return true

}
