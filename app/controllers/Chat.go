//go:build wasm

// Package controllers contains the Gooey controllers of the Exocomp App. Each
// controller drives one View and owns the interaction between the View's
// Component Graph tree and the engine.Session.
package controllers

import app_components "exocomp-app/components"
import engine "exocomp-app/engine"
import app_types "exocomp-app/types"
import "exocomp-app/schemas"
import app_chat "exocomp-app/utils/chat"
import app_cli "exocomp-app/utils/cli"
import app_fmt "exocomp-app/utils/fmt"
import "exocomp-app/views"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/app"
import "github.com/cookiengineer/gooey/components/content"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/layout"
import "github.com/cookiengineer/gooey/components/ui"
import "encoding/json"
import "slices"
import "sort"
import "strconv"
import "strings"
import "syscall/js"
import "time"

// Chat is the Controller of the Chat View.
type Chat struct {
	Main    *app.Main
	Session *engine.Session
	Config  *app_types.Config
	View    *views.Chat
	Role    string

	paused        bool
	timer_label   time.Duration
	timer_session time.Duration

	questions      []*app_types.Question
	question_index int

	tool_call *app_components.ToolCall

	headerActionListener    *components.EventListener
	dialogActionListener    *components.EventListener
	dialogCancelListener    *dom.EventListener
	answerClickListener     *dom.EventListener
	answerCancelListener    *dom.EventListener
	answerDialog            *dom.Element
	promptActionListener    *components.EventListener
	toolCallSuggestListener *components.EventListener
	agentListChangeListener *components.EventListener
}

func NewChat(main *app.Main, view interfaces.View, config *app_types.Config) *Chat {

	app_view, ok := view.(*views.Chat)

	if ok == false || app_view == nil {
		return nil
	}

	session := engine.NewSession(config, main.Client)

	chat := Chat{
		Main:    main,
		Session: session,
		Config:  config,
		View:    app_view,
		Role:    "user",
	}

	return &chat

}

func (chat *Chat) Name() string {
	return "chat"
}

func (chat *Chat) Enter() bool {

	chat.attachHeader()
	chat.attachDialog()
	chat.attachAnswerQuestions()
	chat.attachNavigation()
	chat.attachPopover()
	chat.attachPrompt()

	go chat.Init()

	return true

}

func (chat *Chat) Leave() bool {

	if chat.Main.Header != nil && chat.headerActionListener != nil {
		chat.Main.Header.Component.RemoveEventListener("action", chat.headerActionListener)
		chat.headerActionListener = nil
	}

	if chat.Main.Dialog != nil {

		if chat.dialogActionListener != nil {
			chat.Main.Dialog.Component.RemoveEventListener("action", chat.dialogActionListener)
			chat.dialogActionListener = nil
		}

		if chat.Main.Dialog.Component.Element != nil && chat.dialogCancelListener != nil {
			chat.Main.Dialog.Component.Element.RemoveEventListener("cancel", chat.dialogCancelListener)
			chat.dialogCancelListener = nil
		}

	}

	if chat.answerDialog != nil {

		if chat.answerClickListener != nil {
			chat.answerDialog.RemoveEventListener("click", chat.answerClickListener)
			chat.answerClickListener = nil
		}

		if chat.answerCancelListener != nil {
			chat.answerDialog.RemoveEventListener("cancel", chat.answerCancelListener)
			chat.answerCancelListener = nil
		}

		chat.answerDialog = nil

	}

	if chat.promptActionListener != nil {
		if prompt := chat.prompt(); prompt != nil {
			prompt.Component.RemoveEventListener("action", chat.promptActionListener)
		}
		chat.promptActionListener = nil
	}

	if chat.tool_call != nil && chat.toolCallSuggestListener != nil {
		chat.tool_call.Component.RemoveEventListener("suggest", chat.toolCallSuggestListener)
		chat.toolCallSuggestListener = nil
	}

	if chat.agentListChangeListener != nil {
		if agent_list := chat.View.GetAgentList(); agent_list != nil {
			agent_list.Component.RemoveEventListener("change-agent", chat.agentListChangeListener)
		}
		chat.agentListChangeListener = nil
	}

	return true

}

func (chat *Chat) Render() {
	chat.RenderAll()
}

func (chat *Chat) Update() {
	chat.Session.UpdateAgents()
}

func (chat *Chat) Init() {

	go chat.Session.UpdateTools()
	go chat.Session.UpdateAgents()

	time.Sleep(500 * time.Millisecond)

	chat.Session.UpdateAgents()

	if chat.tool_call != nil {
		chat.tool_call.SetTools(chat.Session.Tools)
	}

	chat.RenderAll()

	go chat.loop()

}

func (chat *Chat) loop() {

	for true {

		time.Sleep(250 * time.Millisecond)

		chat.timer_label += 250 * time.Millisecond
		chat.timer_session += 250 * time.Millisecond

		if chat.paused == true {
			continue
		}

		if chat.timer_label >= 1*time.Second {
			chat.UpdateLabel()
			chat.UpdateQuestions()
			chat.timer_label = 0
		}

		if chat.timer_session >= 5*time.Second {
			chat.Session.UpdateAgents()
			chat.RenderAll()
			chat.timer_session = 0
		}

	}

}

func (chat *Chat) RenderAll() {

	agent := chat.Session.GetAgent("")

	if agent != nil {

		temperature := strconv.FormatFloat(agent.Temperature, 'f', 1, 64)
		header := agent.Name + " | " + agent.Role + " | " + agent.Model + " | " + temperature

		chat.renderHeader(header)
		chat.renderTitle("Exocomp - " + agent.Name)

		if agent.Name == chat.Config.Name {
			chat.enablePrompt()
		} else {
			chat.disablePrompt()
		}

	}

	if agent_list := chat.View.GetAgentList(); agent_list != nil {
		agent_list.SetAgents(chat.Session.Agent, chat.Session.GetAgents())
	}

	if message_list := chat.View.GetMessageList(); message_list != nil {
		message_list.SetMessages(chat.Session.GetMessages(0), chat.Config.Debug)
	}

	chat.UpdateLabel()

}

func (chat *Chat) Pause() {
	chat.paused = true
}

func (chat *Chat) Resume() {
	chat.paused = false
}

func (chat *Chat) SetRole(role string) {

	if role == "assistant" || role == "user" {
		chat.Role = role
	}

}

func (chat *Chat) UpdateLabel() {

	tokens := 0
	length := 0
	cost := 0.0

	agent := chat.Session.GetAgent("")

	if agent != nil {
		tokens = agent.ContextUsage.Tokens
		length = agent.ContextUsage.Length
		cost = agent.ContextUsage.Cost
	}

	value := "0%"
	title := "Context Usage: " + formatTokens(tokens) + " of " + formatTokens(length)

	if length > 0 {
		value = strconv.Itoa(int((float64(tokens)/float64(length))*100)) + "%"
	}

	if cost > 0 {
		title = title + " | Cost: $" + formatCost(cost)
	}

	if chat.Session.Waiting == true {
		value = "\u2026"
		title = "Thinking ..."
	}

	if prompt := chat.prompt(); prompt != nil {
		prompt.SetUsage(value, title)
	}

}

func (chat *Chat) UpdateQuestions() {

	agent := chat.Session.GetAgent("")
	questions := app_chat.GetUnansweredQuestions(agent)

	if len(questions) > 0 {
		chat.ShowQuestions(questions)
		chat.Pause()
	}

}

func (chat *Chat) renderHeader(text string) {

	if chat.Main.Header != nil && chat.Main.Header.Component.Element != nil {
		label := chat.Main.Header.Component.Element.QuerySelector("label[data-name=\"title\"]")
		if label != nil {
			label.SetTextContent(text)
		}
	}

}

func (chat *Chat) renderTitle(text string) {
	js.Global().Get("document").Set("title", text)
}

func (chat *Chat) attachHeader() {

	if chat.Main.Header == nil {
		return
	}

	chat.headerActionListener = components.ToEventListener(func(event string, attributes map[string]any) {

		action, ok := attributes["action"].(string)

		if ok == true && action == "hire-agent" {
			chat.ShowHireDialog()
		}

	}, false)

	chat.Main.Header.Component.AddEventListener("action", chat.headerActionListener)

}

func (chat *Chat) attachDialog() {

	if chat.Main.Dialog == nil {
		return
	}

	chat.dialogActionListener = components.ToEventListener(func(event string, attributes map[string]any) {

		action, ok := attributes["action"].(string)

		if ok == false {
			return
		}

		if action == "confirm" {
			chat.ConfirmHireDialog()
		} else if action == "cancel" || action == "close" {
			chat.resetHireDialog()
			chat.Main.Dialog.Hide()
		}

	}, false)

	chat.Main.Dialog.Component.AddEventListener("action", chat.dialogActionListener)

	if chat.Main.Dialog.Component.Element != nil {

		chat.dialogCancelListener = dom.ToEventListener(func(event *dom.Event) {
			chat.resetHireDialog()
			chat.Main.Dialog.Hide()
		})

		chat.Main.Dialog.Component.Element.AddEventListener("cancel", chat.dialogCancelListener)

	}

}

func (chat *Chat) attachAnswerQuestions() {

	dialog := chat.View.QuerySelector("dialog[data-name=\"answer-questions\"]")

	if dialog == nil {
		return
	}

	chat.answerDialog = dialog

	chat.answerClickListener = dom.ToEventListener(func(event *dom.Event) {

		if event.Target == nil {
			return
		}

		action := event.Target.GetAttribute("data-action")

		if action == "cancel" || action == "close" {
			chat.hideAnswerQuestions()
			chat.Resume()
		} else if action == "confirm" {
			chat.ConfirmQuestion()
		}

	})

	dialog.AddEventListener("click", chat.answerClickListener)

	chat.answerCancelListener = dom.ToEventListener(func(event *dom.Event) {
		chat.hideAnswerQuestions()
		chat.Resume()
	})

	dialog.AddEventListener("cancel", chat.answerCancelListener)

}

func (chat *Chat) attachNavigation() {

	agent_list := chat.View.GetAgentList()

	if agent_list == nil {
		return
	}

	chat.agentListChangeListener = components.ToEventListener(func(event string, attributes map[string]any) {

		name, ok := attributes["name"].(string)

		if ok == true && name != "" {
			chat.ViewAgent(name)
		}

	}, false)

	agent_list.Component.AddEventListener("change-agent", chat.agentListChangeListener)

}

func (chat *Chat) attachPrompt() {

	prompt := chat.prompt()

	if prompt == nil {
		return
	}

	chat.promptActionListener = components.ToEventListener(func(event string, attributes map[string]any) {

		action, _ := attributes["action"].(string)
		value, _ := attributes["value"].(string)

		suggest := true

		if tmp, ok := attributes["suggest"].(bool); ok == true {
			suggest = tmp
		}

		if action == "change" {

			if chat.tool_call != nil {
				chat.tool_call.SetPrompt(value, suggest)
			}

		} else if action == "run-tool" {

			if chat.tool_call != nil {
				chat.tool_call.SetPrompt("", false)
			}

			go chat.runToolCommand(value)

		} else if action == "send" {

			if chat.tool_call != nil {
				chat.tool_call.SetPrompt("", false)
			}

			go chat.sendMessage(value)

		}

	}, false)

	prompt.Component.AddEventListener("action", chat.promptActionListener)

}

func (chat *Chat) attachPopover() {

	component, ok := components.UnwrapComponent[*app_components.ToolCall](
		chat.Main.Document.QueryComponent("body > tool-call"),
	)

	if ok == false || component == nil {
		return
	}

	chat.tool_call = component
	chat.tool_call.SetTools(chat.Session.Tools)

	chat.toolCallSuggestListener = components.ToEventListener(func(event string, attributes map[string]any) {

		if prompt := chat.prompt(); prompt != nil {
			prompt.ApplySuggestion(attributes)
		}

	}, false)

	chat.tool_call.Component.AddEventListener("suggest", chat.toolCallSuggestListener)

}

func (chat *Chat) sendMessage(prompt string) {

	ok := chat.Session.SendChatRequest(&schemas.Message{
		Role:    chat.Role,
		Content: prompt,
	})

	if ok == true {
		chat.Session.UpdateAgents()
		chat.RenderAll()
	}

}

func (chat *Chat) runToolCommand(prompt string) {

	trimmed := strings.TrimPrefix(prompt, "/")

	var name string
	var method string
	var args map[string]any

	if strings.Contains(trimmed, " ") == true {

		name = strings.Split(trimmed, " ")[0]
		args = app_cli.ParseParameters(strings.TrimPrefix(trimmed, name))

	} else {

		name = strings.TrimSpace(trimmed)
		args = make(map[string]any)

	}

	if strings.Contains(name, ".") == true {
		method = name[strings.LastIndex(name, ".")+1:]
	}

	if name == "" || method == "" {
		return
	}

	chat.Session.CallTool(name, method, args)

}

func (chat *Chat) UpdatePrompt(message string) {

	prompt_value := strings.TrimSpace(message)

	if prompt := chat.prompt(); prompt != nil {
		prompt.SetValue(prompt_value)
	}

	if chat.tool_call != nil {
		chat.tool_call.SetPrompt(prompt_value, false)
	}

	chat.UpdateLabel()

}

func (chat *Chat) ViewAgent(name string) {

	active := chat.Session.GetAgent("")
	agent := chat.Session.GetAgent(name)

	if agent == nil || agent == active {
		return
	}

	chat.Session.SetAgent(agent.Name)

	temperature := strconv.FormatFloat(agent.Temperature, 'f', 1, 64)

	chat.renderHeader(agent.Name + " | " + agent.Role + " | " + agent.Model + " | " + temperature)

	if agent_list := chat.View.GetAgentList(); agent_list != nil {
		agent_list.SetAgents(chat.Session.Agent, chat.Session.GetAgents())
	}

	if message_list := chat.View.GetMessageList(); message_list != nil {
		message_list.SetMessages(chat.Session.GetMessages(0), chat.Config.Debug)
	}

	if agent.Name == chat.Config.Name {
		chat.enablePrompt()
	} else {
		chat.disablePrompt()
	}

}

func (chat *Chat) HireAgent(data map[string]string) (bool, []string) {

	name := strings.TrimSpace(data["name"])
	role := strings.TrimSpace(data["role"])
	sandbox := strings.TrimSpace(data["sandbox"])
	prompt := strings.TrimSpace(data["prompt"])

	errors := make([]string, 0)

	if name == "" || chat.Session.GetAgent(name) != nil {
		errors = append(errors, "Invalid Agent Name, must be a unique Pseudonym.")
	}

	if role == "" {
		errors = append(errors, "Invalid Agent Role.")
	}

	if strings.HasPrefix(sandbox, "./") == false {
		errors = append(errors, "Invalid Agent Sandbox, must start with \"./path/to/sandbox\"")
	}

	if prompt == "" {
		errors = append(errors, "Invalid Agent Prompt, must not be empty.")
	}

	if len(errors) > 0 {
		return false, errors
	}

	go chat.Session.CallTool("agents.Hire", "Hire", map[string]any{
		"name":    name,
		"role":    role,
		"sandbox": sandbox,
		"prompt":  prompt,
	})

	return true, errors

}

func (chat *Chat) AnswerQuestion(question string, answer string) (bool, []string) {

	question = strings.TrimSpace(question)
	answer = strings.TrimSpace(answer)

	if question == "" || answer == "" {
		return false, []string{"Invalid Answer, must not be empty."}
	}

	go chat.Session.CallTool("humans.Answer", "Answer", map[string]any{
		"question": question,
		"answer":   answer,
	})

	return true, []string{}

}

func (chat *Chat) ShowHireDialog() {

	chat.populateRoles()
	chat.clearHireErrors()

	if chat.Main.Dialog != nil {
		chat.Main.Dialog.Show()
	}

}

func (chat *Chat) populateRoles() {

	if chat.Main.Dialog == nil || chat.Main.Dialog.Component.Element == nil {
		return
	}

	go func() {

		roles := []string{"architect", "coder", "tester"}

		response, err := chat.Main.Client.Read("/api/parameters/roles")

		if err == nil && response != nil {

			fetched := make([]string, 0)

			if json.Unmarshal(response.Body, &fetched) == nil {

				for _, role := range fetched {

					if role == "planner" {
						continue
					}

					if slices.Contains(roles, role) == false {
						roles = append(roles, role)
					}

				}

			}

		}

		sort.Strings(roles)

		select_component, ok := components.UnwrapComponent[*ui.Select](
			chat.Main.Dialog.Query("dialog > fieldset > select"),
		)

		if ok == true && select_component != nil {
			select_component.SetValues(roles)
			select_component.Render()
		}

	}()

}

func (chat *Chat) ConfirmHireDialog() {

	fieldset, ok := components.UnwrapComponent[*content.Fieldset](
		chat.Main.Dialog.Query("dialog > fieldset"),
	)

	if ok == false || fieldset == nil {
		return
	}

	name := fieldValue(fieldset, "name")
	role := fieldValue(fieldset, "role")
	sandbox := fieldValue(fieldset, "sandbox")
	prompt := fieldValue(fieldset, "prompt")

	result, errors := chat.HireAgent(map[string]string{
		"name":    name,
		"role":    role,
		"sandbox": sandbox,
		"prompt":  prompt,
	})

	if result == true {

		fieldset.Reset()
		chat.clearHireErrors()
		chat.Main.Dialog.Hide()
		chat.Session.UpdateAgents()
		chat.RenderAll()

	} else {
		chat.showHireErrors(errors)
	}

}

func (chat *Chat) resetHireDialog() {

	fieldset, ok := components.UnwrapComponent[*content.Fieldset](
		chat.Main.Dialog.Query("dialog > fieldset"),
	)

	if ok == true && fieldset != nil {
		fieldset.Reset()
	}

	chat.clearHireErrors()

}

func (chat *Chat) clearHireErrors() {

	if chat.Main.Dialog == nil {
		return
	}

	clearFooterErrors(chat.Main.Dialog.Footer)

}

func (chat *Chat) showHireErrors(errors []string) {

	if chat.Main.Dialog == nil {
		return
	}

	setFooterErrors(chat.Main.Dialog.Footer, errors)

}

// questionDialog returns the answer-questions dialog of the Chat view.
func (chat *Chat) questionDialog() *layout.Dialog {

	dialog, ok := components.UnwrapComponent[*layout.Dialog](
		chat.View.Query("section > dialog"),
	)

	if ok == true && dialog != nil {
		return dialog
	}

	return nil

}

func (chat *Chat) ShowQuestions(questions []*app_types.Question) {

	chat.questions = questions
	chat.question_index = 0
	chat.RenderQuestion()

}

func (chat *Chat) RenderQuestion() {

	dialog := chat.View.QuerySelector("dialog[data-name=\"answer-questions\"]")

	if dialog == nil || len(chat.questions) == 0 {
		return
	}

	if chat.question_index < 0 {
		chat.question_index = 0
	} else if chat.question_index >= len(chat.questions) {
		chat.question_index = len(chat.questions) - 1
	}

	question := chat.questions[chat.question_index]

	title := dialog.QuerySelector("h3")
	label := dialog.QuerySelector("p[data-name=\"question\"]")
	answers := dialog.QuerySelector("div[data-name=\"answers\"]")

	if title != nil {
		if question.IsChoice() == true {
			title.SetTextContent("Choices")
		} else {
			title.SetTextContent("Question")
		}
	}

	if label != nil {
		label.SetInnerHTML(app_fmt.SanitizeContent(question.Question))
	}

	chat.clearQuestionErrors()

	if answers != nil {

		html := ""

		if question.IsChoice() == true {

			input_type := "radio"

			if question.Multiple == true {
				input_type = "checkbox"
			}

			for _, option := range question.Options {
				html += "<label><input type=\"" + input_type + "\" name=\"question-answer\" value=\"" + app_fmt.SanitizeContent(option) + "\"/><span>" + app_fmt.SanitizeContent(option) + "</span></label>"
			}

		} else {
			html += "<textarea data-name=\"answer\" placeholder=\"Your answer ...\"></textarea>"
		}

		answers.SetInnerHTML(html)

	}

	dialog.Value.Call("showModal")

}

func (chat *Chat) ConfirmQuestion() {

	dialog := chat.View.QuerySelector("dialog[data-name=\"answer-questions\"]")

	if dialog == nil || len(chat.questions) == 0 {
		return
	}

	question := chat.questions[chat.question_index]

	answer := ""

	if question.IsChoice() == true {

		checked := dialog.QuerySelectorAll("input[name=\"question-answer\"]:checked")
		values := make([]string, 0)

		for _, input := range checked {
			values = append(values, input.Value.Get("value").String())
		}

		answer = strings.Join(values, "\n")

	} else {
		answer = valueOf(dialog, "div[data-name=\"answers\"] textarea")
	}

	result, errors := chat.AnswerQuestion(question.Question, answer)

	if result == false {
		chat.showQuestionErrors(errors)
		return
	}

	if chat.question_index >= len(chat.questions)-1 {
		chat.hideAnswerQuestions()
		chat.Resume()
	} else {
		chat.question_index++
		chat.RenderQuestion()
	}

}

func (chat *Chat) hideAnswerQuestions() {

	dialog := chat.View.QuerySelector("dialog[data-name=\"answer-questions\"]")

	if dialog != nil {
		dialog.Value.Call("close")
	}

}

func (chat *Chat) clearQuestionErrors() {

	dialog := chat.questionDialog()

	if dialog != nil {
		clearFooterErrors(dialog.Footer)
	}

}

func (chat *Chat) showQuestionErrors(errors []string) {

	dialog := chat.questionDialog()

	if dialog != nil {
		setFooterErrors(dialog.Footer, errors)
	}

}

// setFooterErrors writes validation errors into the mapped error label of a
// dialog footer. The label is part of the footer Component Graph, so it
// survives the dialog's reconciliation.
func setFooterErrors(footer *layout.Footer, errors []string) {

	if footer == nil {
		return
	}

	label, ok := components.UnwrapComponent[*ui.Label](footer.Query("footer > label"))

	if ok == false || label == nil || label.Component.Element == nil {
		return
	}

	html := ""

	for _, err := range errors {
		html += "<b>" + err + "</b>"
	}

	label.Label = html
	label.Component.Element.SetInnerHTML(html)

}

func clearFooterErrors(footer *layout.Footer) {

	if footer == nil {
		return
	}

	label, ok := components.UnwrapComponent[*ui.Label](footer.Query("footer > label"))

	if ok == false || label == nil || label.Component.Element == nil {
		return
	}

	label.Label = ""
	label.Component.Element.SetInnerHTML("")

}

// prompt returns the ChatPrompt mapped inside the Chat view's Component Graph.
func (chat *Chat) prompt() *app_components.ChatPrompt {

	prompt, ok := components.UnwrapComponent[*app_components.ChatPrompt](
		chat.View.Query("section > chat-footer > chat-prompt"),
	)

	if ok == true && prompt != nil {
		return prompt
	}

	return nil

}

func (chat *Chat) enablePrompt() {

	if prompt := chat.prompt(); prompt != nil {
		prompt.Enable()
	}

}

func (chat *Chat) disablePrompt() {

	if prompt := chat.prompt(); prompt != nil {
		prompt.Disable()
	}

}

func fieldValue(fieldset *content.Fieldset, name string) string {

	value := fieldset.ValueOf(name)

	if value.IsUndefined() || value.IsNull() {
		return ""
	}

	return strings.TrimSpace(value.String())

}

func valueOf(root *dom.Element, query string) string {

	element := root.QuerySelector(query)

	if element == nil || element.Value == nil {
		return ""
	}

	return strings.TrimSpace(element.Value.Get("value").String())

}

func formatCost(cost float64) string {

	value := strconv.FormatFloat(cost, 'f', 4, 64)

	if strings.Contains(value, ".") == true {

		value = strings.TrimRight(value, "0")
		value = strings.TrimRight(value, ".")

	}

	return value

}

// formatTokens renders a token count in a compact human readable form, e.g.
// 123 -> "123", 123456 -> "123k", 1234567 -> "1.2M".
func formatTokens(tokens int) string {

	if tokens < 0 {
		tokens = 0
	}

	if tokens >= 1000000 {

		value := strconv.FormatFloat(float64(tokens)/1000000, 'f', 1, 64)
		value = strings.TrimRight(value, "0")
		value = strings.TrimRight(value, ".")

		return value + "M"

	} else if tokens >= 1000 {

		return strconv.Itoa(tokens/1000) + "k"

	}

	return strconv.Itoa(tokens)

}
