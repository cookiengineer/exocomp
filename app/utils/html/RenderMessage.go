package html

import "exocomp-app/schemas"
import "encoding/json"
import std_html "html"
import "sort"
import "strings"

func RenderMessage(message *schemas.Message, with_empty_content bool) string {

	if message == nil {
		return ""
	}

	if message.Role == "assistant" {

		html := ""

		if message.ReasoningContent != "" {
			html += renderCollapse("reasoning", "Thought", previewText(message.ReasoningContent), RenderMarkdown(message.ReasoningContent))
		}

		if message.Content != "" {
			html += RenderMarkdown(message.Content)
		}

		for _, tool_call := range message.ToolCalls {
			html += renderToolCall(tool_call)
		}

		if html != "" {
			return html
		} else if with_empty_content == true {
			return "(no content)"
		}

	} else if message.Role == "system" {

		if message.Content != "" {
			return renderCollapse("system", "System", previewText(message.Content), RenderMarkdown(message.Content))
		} else if with_empty_content == true {
			return "(no content)"
		}

	} else if message.Role == "tool" {

		label := "Tool"

		if strings.TrimSpace(message.ToolName) != "" {
			label = "Tool · " + strings.TrimSpace(message.ToolName)
		}

		if message.Content != "" {
			return renderCollapse("tool-output", label, previewText(message.Content), "<pre>"+std_html.EscapeString(message.Content)+"</pre>")
		}

		return renderCollapse("tool-output", label, "(no content)", "<pre>(no content)</pre>")

	} else if message.Role == "user" {

		if message.Content != "" {
			return RenderMarkdown(message.Content)
		} else if with_empty_content == true {
			return "(no content)"
		}

	}

	return ""

}

func renderCollapse(kind string, label string, preview string, body string) string {

	summary := "<summary><span class=\"type\">" + std_html.EscapeString(label) + "</span>"

	if preview != "" {
		summary += "<span class=\"preview\">" + std_html.EscapeString(preview) + "</span>"
	}

	summary += "</summary>"

	return "<details class=\"message-collapse " + kind + "\">" + summary + body + "</details>"

}

func renderToolCall(tool_call schemas.ToolCall) string {

	name := strings.TrimSpace(tool_call.Function.Name)

	if name == "" {
		name = "Tool"
	}

	return renderCollapse(
		"tool-call",
		"Tool · "+name,
		toolCallPreview(tool_call),
		"<pre>"+std_html.EscapeString(toolCallContent(tool_call))+"</pre>",
	)

}

func previewText(content string) string {

	lines := strings.Split(content, "\n")
	preview := ""
	remaining := false

	for _, line := range lines {

		line = strings.Join(strings.Fields(line), " ")

		if line == "" {
			continue
		}

		if preview == "" {
			preview = line
		} else {
			remaining = true
			break
		}

	}

	if remaining == true {
		return preview + "…"
	}

	return preview

}

func toolCallPreview(tool_call schemas.ToolCall) string {

	arguments, err := tool_call.GetArguments()

	if err == nil {

		keys := make([]string, 0, len(arguments))

		for key := range arguments {
			keys = append(keys, key)
		}

		if len(keys) > 0 {

			sort.Strings(keys)
			key := keys[0]
			value := fmtValue(arguments[key])

			return key + "=" + strings.Join(strings.Fields(value), " ")

		}

	}

	return strings.Join(strings.Fields(string(tool_call.Function.ArgumentsRaw)), " ")

}

func toolCallContent(tool_call schemas.ToolCall) string {

	arguments, err := tool_call.GetArguments()

	if err == nil {

		bytes, err0 := json.MarshalIndent(arguments, "", "  ")

		if err0 == nil {
			return string(bytes)
		}

	}

	return strings.TrimSpace(string(tool_call.Function.ArgumentsRaw))

}

func fmtValue(value any) string {

	if value == nil {
		return "null"
	}

	switch value.(type) {
	case string:
		return value.(string)
	default:
		bytes, err := json.Marshal(value)

		if err == nil {
			return string(bytes)
		}

	}

	return ""

}
