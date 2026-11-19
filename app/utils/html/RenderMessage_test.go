package html

import "exocomp-app/schemas"
import "strings"
import "testing"

func TestRenderMessageSystemCollapsed(t *testing.T) {

	message := &schemas.Message{
		Role:    "system",
		Content: "You are an Exocomp robot.\nThis is the full system prompt.",
	}

	html := RenderMessage(message, false)

	if strings.Contains(html, "message-collapse system") == false {
		t.Errorf("Expected a collapsed system block, got %q", html)
	}

	if strings.Contains(html, ">System<") == false {
		t.Errorf("Expected the System type label, got %q", html)
	}

	if strings.Contains(html, "You are an Exocomp robot.…") == false {
		t.Errorf("Expected the first line preview, got %q", html)
	}

	if strings.Contains(html, "This is the full system prompt.") == false {
		t.Errorf("Expected the full system content, got %q", html)
	}

}

func TestRenderMessageToolCollapsed(t *testing.T) {

	message := &schemas.Message{
		Role:     "tool",
		ToolName: "files.Read",
		Content:  "first line\nsecond line",
	}

	html := RenderMessage(message, false)

	if strings.Contains(html, "message-collapse tool-output") == false {
		t.Errorf("Expected a collapsed tool block, got %q", html)
	}

	if strings.Contains(html, "Tool · files.Read") == false {
		t.Errorf("Expected the tool type, got %q", html)
	}

	if strings.Contains(html, "first line…") == false {
		t.Errorf("Expected the first line preview, got %q", html)
	}

	if strings.Contains(html, "second line") == false {
		t.Errorf("Expected the full tool content, got %q", html)
	}

}

func TestRenderMessageToolEscapesContent(t *testing.T) {

	message := &schemas.Message{
		Role:     "tool",
		ToolName: "files.Read",
		Content:  "<unsafe>",
	}

	html := RenderMessage(message, false)

	if strings.Contains(html, "&lt;unsafe&gt;") == false {
		t.Errorf("Expected escaped tool content, got %q", html)
	}

	if strings.Contains(html, "<unsafe>") == true {
		t.Errorf("Expected no raw HTML in tool content, got %q", html)
	}

}

func TestRenderMessageAssistantReasoningAndToolCall(t *testing.T) {

	message := &schemas.Message{
		Role:             "assistant",
		Content:          "The answer.",
		ReasoningContent: "I should inspect the files first.",
		ToolCalls: []schemas.ToolCall{
			{
				ID:   "call_0",
				Type: "function",
				Function: schemas.ToolCallFunction{
					Name:         "files.Read",
					ArgumentsRaw: []byte(`{"path":"/tmp/file.go"}`),
				},
			},
		},
	}

	html := RenderMessage(message, false)

	if strings.Contains(html, "message-collapse reasoning") == false {
		t.Errorf("Expected a collapsed reasoning block, got %q", html)
	}

	if strings.Contains(html, ">Thought<") == false {
		t.Errorf("Expected the Thought label, got %q", html)
	}

	if strings.Contains(html, "message-collapse tool-call") == false {
		t.Errorf("Expected a collapsed tool call, got %q", html)
	}

	if strings.Contains(html, "files.Read") == false {
		t.Errorf("Expected the tool call name, got %q", html)
	}

	if strings.Contains(html, "/tmp/file.go") == false {
		t.Errorf("Expected the tool call arguments, got %q", html)
	}

}
