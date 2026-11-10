//go:build wasm

package components

import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import "fmt"
import std_html "html"
import "strconv"
import "strings"
import "syscall/js"

// ChatPrompt is the custom prompt component of the Chat view. It owns the
// `<textarea data-name="prompt">` and the `<label data-name="usage">` and
// translates the prompt keystrokes into `action` events, so that the Chat
// controller only has to react to `change`, `run-tool` and `send`.
type ChatPrompt struct {
	Component *components.Component  `json:"component"`
	Content   []interfaces.Component `json:"content"`
	Value     string                 `json:"value"`
	Usage     string                 `json:"usage"`
	Command   bool                   `json:"command"`
	Disabled  bool                   `json:"disabled"`

	label    *dom.Element `json:"-"`
	textarea *dom.Element `json:"-"`

	on_keyup *dom.EventListener `json:"-"`
	on_input *dom.EventListener `json:"-"`
}

func NewChatPrompt() ChatPrompt {

	element := dom.GetDocument().CreateElement("chat-prompt")
	component := components.NewComponent(element)

	self := ChatPrompt{
		Component: &component,
		Content:   make([]interfaces.Component, 0),
	}

	return self

}

func ToChatPrompt(element *dom.Element) *ChatPrompt {

	var self ChatPrompt

	component := components.NewComponent(element)

	self.Component = &component
	self.Content = make([]interfaces.Component, 0)

	return &self

}

func (prompt *ChatPrompt) Disable() bool {

	prompt.Disabled = true

	if prompt.textarea != nil {
		prompt.textarea.SetAttribute("disabled", "")
	}

	return true

}

func (prompt *ChatPrompt) Enable() bool {

	prompt.Disabled = false

	if prompt.textarea != nil {
		prompt.textarea.RemoveAttribute("disabled")
	}

	return true

}

func (prompt *ChatPrompt) Invalidate() {
	prompt.Component.InvalidateAs(prompt)
}

func (prompt *ChatPrompt) SetScheduler(scheduler interfaces.Scheduler) {

	if prompt.Component != nil {
		prompt.Component.SetScheduler(scheduler)
	}

	for _, content := range prompt.Content {
		content.SetScheduler(scheduler)
	}

}

func (prompt *ChatPrompt) Mount() bool {

	if prompt.Component == nil || prompt.Component.Element == nil {
		return false
	}

	prompt.Component.InitEvent("action")

	elements := prompt.Component.Element.Children()
	content := make([]interfaces.Component, 0)

	for _, element := range elements {

		switch element.TagName {

		case "LABEL":
			prompt.label = element

		case "TEXTAREA":
			prompt.textarea = element

		}

		tmp := components.NewComponent(element)
		content = append(content, &tmp)

	}

	prompt.Content = content

	for _, child := range prompt.Content {
		child.Mount()
	}

	if prompt.textarea != nil {

		prompt.Disabled = prompt.textarea.HasAttribute("disabled")

		if value := prompt.textarea.Value.Get("value"); !value.IsNull() && !value.IsUndefined() {
			prompt.Value = value.String()
		}

		prompt.on_keyup = dom.ToEventListener(func(event *dom.Event) {
			prompt.HandleKeyup(event)
		})

		prompt.textarea.AddEventListener("keyup", prompt.on_keyup)

		prompt.on_input = dom.ToEventListener(func(event *dom.Event) {
			prompt.Resize()
		})

		prompt.textarea.AddEventListener("input", prompt.on_input)

	}

	if prompt.label != nil {
		prompt.Usage = strings.TrimSpace(prompt.label.TextContent)
	}

	prompt.Resize()

	return true

}

// HandleKeyup mirrors the old frontend `Client.mjs` prompt `keyup` handler and
// emits `change`, `run-tool` or `send` actions.
func (prompt *ChatPrompt) HandleKeyup(event *dom.Event) {

	if event == nil || event.Value.IsNull() || event.Value.IsUndefined() {
		return
	}

	key := event.Value.Get("key").String()
	ctrl := event.Value.Get("ctrlKey").Bool()

	value := prompt.GetValue()

	prompt.Value = value
	prompt.Command = strings.HasPrefix(value, "/") == true

	prompt.applyCommandState()

	if key == "Enter" {

		has_space := strings.Contains(value, " ")
		has_newline := strings.Contains(value, "\n")

		if value == "" {
			return
		}

		if ctrl == true {

			if strings.HasPrefix(value, "/") == true && has_newline == false {

				prompt.Component.FireEventListeners("action", map[string]any{
					"action": "run-tool",
					"value":  value,
				})

				prompt.SetValue("")

			} else {

				prompt.Component.FireEventListeners("action", map[string]any{
					"action": "send",
					"value":  value,
				})

				prompt.SetValue("")

			}

		} else if strings.HasPrefix(value, "/") == true && has_space == true && has_newline == false {

			prompt.Component.FireEventListeners("action", map[string]any{
				"action": "run-tool",
				"value":  value,
			})

			prompt.SetValue("")

		}

		return

	}

	if key == "Backspace" || key == "Delete" {

		// NOTE: Re-render the suggestions, but don't auto-complete the
		// current parameter while the user is deleting.
		prompt.Component.FireEventListeners("action", map[string]any{
			"action":  "change",
			"value":   value,
			"suggest": false,
		})

		return

	}

	prompt.Component.FireEventListeners("action", map[string]any{
		"action":  "change",
		"value":   value,
		"suggest": true,
	})

}

func (prompt *ChatPrompt) GetValue() string {

	if prompt.textarea != nil {

		value := prompt.textarea.Value.Get("value")

		if !value.IsNull() && !value.IsUndefined() {
			return strings.TrimSpace(value.String())
		}

	}

	return prompt.Value

}

func (prompt *ChatPrompt) SetValue(value string) {

	prompt.Value = value
	prompt.Command = strings.HasPrefix(value, "/") == true

	if prompt.textarea != nil {
		prompt.textarea.Value.Set("value", value)
	}

	prompt.applyCommandState()
	prompt.Resize()
	prompt.Invalidate()

}

// Resize grows or shrinks the textarea so that it always fits its content,
// starting at a single line and expanding on wrapped text or newlines.
func (prompt *ChatPrompt) Resize() {

	if prompt.textarea == nil || prompt.textarea.Value == nil {
		return
	}

	style := prompt.textarea.Value.Get("style")

	style.Set("height", "auto")

	scroll_height := prompt.textarea.Value.Get("scrollHeight").Int()
	offset_height := prompt.textarea.Value.Get("offsetHeight").Int()
	client_height := prompt.textarea.Value.Get("clientHeight").Int()

	border := 0

	if offset_height > client_height {
		border = offset_height - client_height
	}

	if scroll_height < 1 {
		scroll_height = 1
	}

	style.Set("height", strconv.Itoa(scroll_height+border)+"px")

	prompt.publishHeight()

}

// publishHeight exposes the rendered ChatFooter height as the
// `--chat-footer-height` CSS variable so that the message list, the agent
// sidebar and the tool popover can stay above the growing footer.
func (prompt *ChatPrompt) publishHeight() {

	if prompt.Component == nil || prompt.Component.Element == nil {
		return
	}

	parent := prompt.Component.Element.ParentNode()

	if parent == nil || parent.Value == nil {
		return
	}

	height := parent.Value.Get("offsetHeight").Int()

	if height < 1 {
		return
	}

	js.Global().Get("document").Get("documentElement").Get("style").Call(
		"setProperty",
		"--chat-footer-height",
		strconv.Itoa(height)+"px",
	)

}

func (prompt *ChatPrompt) SetUsage(value string, title string) {

	prompt.Usage = value

	if prompt.label != nil {

		if title != "" {
			prompt.label.SetInnerHTML(fmt.Sprintf("<abbr title=\"%s\">%s</abbr>", std_html.EscapeString(title), std_html.EscapeString(value)))
		} else {
			prompt.label.SetInnerHTML(std_html.EscapeString(value))
		}

	}

	prompt.Resize()
	prompt.Invalidate()

}

func (prompt *ChatPrompt) Focus() {

	if prompt.textarea != nil {
		prompt.textarea.Value.Call("focus")
	}

}

// ApplySuggestion mirrors the old frontend `Client.mjs` `Suggest` method. It
// completes the current slash-command parameter with the given suggestion.
func (prompt *ChatPrompt) ApplySuggestion(suggestion map[string]any) {

	value := prompt.GetValue()

	if strings.HasPrefix(value, "/") == false || strings.Contains(value, " ") == false {
		return
	}

	parts := strings.Split(value, " ")
	current_parameter := ""

	if len(parts) > 1 {
		current_parameter = parts[len(parts)-1]
	}

	typ, _ := suggestion["type"].(string)
	key, _ := suggestion["key"].(string)
	suggestion_value, _ := suggestion["value"].(string)

	next_suggestion := ""

	if typ == "array-of-booleans" || typ == "array-of-numbers" || typ == "array-of-strings" {

		if strings.HasSuffix(current_parameter, "=") == true {

			next_suggestion = suggestion_value

		} else if strings.Contains(current_parameter, "=") == false {

			if strings.HasPrefix(key, current_parameter) == true && len(current_parameter) < len(key) {
				next_suggestion = key[len(current_parameter):]
			} else if current_parameter == key {
				next_suggestion = "=["
			}

		}

	} else if typ == "boolean" || typ == "number" || typ == "string" {

		if strings.HasSuffix(current_parameter, "=") == true {

			next_suggestion = suggestion_value

		} else if strings.Contains(current_parameter, "=") == false {

			if strings.HasPrefix(key, current_parameter) == true && len(current_parameter) < len(key) {
				next_suggestion = key[len(current_parameter):]
			} else if current_parameter == key {
				next_suggestion = "="
			}

		}

	}

	if next_suggestion != "" && prompt.textarea != nil {

		prompt.textarea.Value.Set("value", value+next_suggestion)
		prompt.textarea.Value.Call("setSelectionRange", len(value), len(value)+len(next_suggestion))
		prompt.textarea.Value.Call("focus")

	}

}

func (prompt *ChatPrompt) applyCommandState() {

	if prompt.Component == nil || prompt.Component.Element == nil {
		return
	}

	if prompt.Command == true {
		prompt.Component.Element.SetAttribute("data-command", "true")
	} else {
		prompt.Component.Element.RemoveAttribute("data-command")
	}

}

func (prompt *ChatPrompt) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) >= 2 {

		if prompt.Component.Element != nil {

			if utils.MatchesQuery(prompt.Component.Element, selectors[0]) == true {

				tmp_query := utils.JoinQuery(selectors[1:])

				for _, content := range prompt.Content {

					tmp_component := content.Query(tmp_query)

					if tmp_component != nil {
						return tmp_component
					}

				}

			}

		}

	} else if len(selectors) == 1 {

		if prompt.Component.Element != nil {

			if utils.MatchesQuery(prompt.Component.Element, selectors[0]) == true {
				return prompt
			}

		}

	}

	return nil

}

func (prompt *ChatPrompt) Render() *dom.Element {

	if prompt.Component.Element != nil {

		prompt.applyCommandState()

		components.ReconcileComponents(prompt.Component.Element, prompt.Content)

	}

	return prompt.Component.Element

}

func (prompt *ChatPrompt) String() string {

	html := "<chat-prompt>"

	for _, content := range prompt.Content {
		html += content.String()
	}

	html += "</chat-prompt>"

	return html

}

func (prompt *ChatPrompt) Unmount() bool {

	if prompt.textarea != nil {

		if prompt.on_keyup != nil {
			prompt.textarea.RemoveEventListener("keyup", prompt.on_keyup)
			prompt.on_keyup = nil
		}

		if prompt.on_input != nil {
			prompt.textarea.RemoveEventListener("input", prompt.on_input)
			prompt.on_input = nil
		}

	}

	for _, content := range prompt.Content {
		content.Unmount()
	}

	return true

}
