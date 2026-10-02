//go:build wasm

package components

import "github.com/cookiengineer/gooey/components"

func RegisterTo(document *components.Document) {

	document.Register("agent-list", components.WrapComponent(ToAgentList))
	document.Register("tool-call", components.WrapComponent(ToToolCall))
	document.Register("message-list", components.WrapComponent(ToMessageList))
	document.Register("chat-footer", components.WrapComponent(ToChatFooter))
	document.Register("chat-prompt", components.WrapComponent(ToChatPrompt))
	document.Register("bug-grid", components.WrapComponent(ToBugGrid))
	document.Register("requirement-grid", components.WrapComponent(ToRequirementGrid))
	document.Register("changelog-grid", components.WrapComponent(ToChangelogGrid))

}
