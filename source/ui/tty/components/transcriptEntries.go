package components

import "exocomp/schemas"
import "exocomp/ui/tty/colors"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "encoding/json"
import "fmt"
import "sort"
import "strings"

func (transcript *Transcript) appendMessage(result []transcriptLine, message *schemas.Message, width int) []transcriptLine {

	switch message.Role {

	case "system":

		result = append(result, withBorder(transcript.collapsibleLines(
			message,
			"system",
			"System",
			firstContentLine(message.Content),
			message.Content,
			styles.SystemHeader,
			styles.CollapsedHint,
			width,
			false,
		), styles.BorderSystem)...)

	case "tool":

		label := "Tool"

		if strings.TrimSpace(message.ToolName) != "" {
			label = "Tool · " + strings.TrimSpace(message.ToolName)
		}

		result = append(result, withBorder(transcript.collapsibleLines(
			message,
			"tool",
			label,
			firstContentLine(message.Content),
			message.Content,
			styles.ToolHeader,
			styles.CollapsedHint,
			width,
			false,
		), styles.BorderTool)...)

	case "assistant":

		result = append(result, transcriptLine{
			Spans: []transcriptSpan{
				{Text: "[assistant]", Style: styles.TranscriptAssistant},
			},
			HasBorder: true,
			Border:    styles.BorderAssistant,
		})

		if strings.TrimSpace(message.ReasoningContent) != "" {

			result = append(result, withBorder(transcript.collapsibleLines(
				message,
				"reasoning",
				"Thought",
				firstContentLine(message.ReasoningContent),
				message.ReasoningContent,
				styles.ReasoningHeader,
				styles.CollapsedHint,
				width,
				true,
			), styles.BorderThinking)...)

		}

		if strings.TrimSpace(message.Content) != "" {
			result = append(result, withBorder(renderMarkdown(message.Content, width), styles.BorderAssistant)...)
		}

		for index, tool_call := range message.ToolCalls {

			name := strings.TrimSpace(tool_call.Function.Name)

			if name == "" {
				name = "Tool"
			}

			result = append(result, withBorder(transcript.collapsibleLines(
				message,
				fmt.Sprintf("toolcall:%d", index),
				"Tool · "+name,
				toolCallPreview(tool_call),
				toolCallContent(tool_call),
				styles.ToolHeader,
				styles.CollapsedHint,
				width,
				false,
			), styles.BorderTool)...)

		}

	case "user":

		result = append(result, transcriptLine{
			Spans: []transcriptSpan{
				{Text: "[user]", Style: styles.TranscriptUser},
			},
			HasBorder: true,
			Border:    styles.BorderUser,
		})

		if strings.TrimSpace(message.Content) != "" {
			result = append(result, withBorder(renderMarkdown(message.Content, width), styles.BorderUser)...)
		}

	default:

		result = append(result, transcriptLine{
			Spans: []transcriptSpan{
				{Text: fmt.Sprintf("[%s]", message.Role), Style: styles.TranscriptTool},
			},
			HasBorder: true,
			Border:    styles.BorderTool,
		})

		if strings.TrimSpace(message.Content) != "" {
			result = append(result, withBorder(renderMarkdown(message.Content, width), styles.BorderTool)...)
		}

	}

	return result

}

func withBorder(lines []transcriptLine, style utils_tty.Style) []transcriptLine {

	for index := range lines {
		lines[index].HasBorder = true
		lines[index].Border = style
	}

	return lines

}

func (transcript *Transcript) collapsibleLines(message *schemas.Message, part string, label string, preview string, content string, header_style utils_tty.Style, hint_style utils_tty.Style, width int, mute bool) []transcriptLine {

	toggle := transcriptToggle{
		Message: message,
		Part:    part,
	}

	if transcript.expanded[toggle] == false {

		spans := []transcriptSpan{
			{Text: "▸ " + label, Style: header_style},
		}

		preview = strings.Join(strings.Fields(preview), " ")

		if preview != "" {

			prefix_width := utils_tty.StringWidth("▸ " + label + ": ")
			preview = utils_tty.TruncateRight(preview, max(1, width-prefix_width))
			spans = append(spans, transcriptSpan{Text: ": " + preview, Style: hint_style})

		}

		return []transcriptLine{
			{
				Spans:  compactSpans(spans),
				Toggle: &toggle,
			},
		}

	}

	header := transcriptLine{
		Spans: []transcriptSpan{
			{Text: "▾ " + label, Style: header_style},
		},
		Toggle: &toggle,
	}

	if strings.TrimSpace(content) == "" {
		return []transcriptLine{header}
	}

	body := renderMarkdown(content, width)

	if mute == true {
		body = muteTranscriptLines(body)
	}

	body = indentTranscriptLines(body, "  ", hint_style)

	return append([]transcriptLine{header}, body...)

}

func firstContentLine(content string) string {

	for _, line := range strings.Split(content, "\n") {

		line = strings.Join(strings.Fields(line), " ")

		if line != "" {
			return line
		}

	}

	return ""

}

func toolCallPreview(tool_call schemas.ToolCall) string {

	raw := strings.TrimSpace(string(tool_call.Function.ArgumentsRaw))

	if raw == "" {
		return ""
	}

	arguments, err := tool_call.GetArguments()

	if err == nil {

		keys := make([]string, 0, len(arguments))

		for key := range arguments {
			keys = append(keys, key)
		}

		if len(keys) > 0 {

			sort.Strings(keys)
			key := keys[0]
			value := fmt.Sprintf("%v", arguments[key])

			return key + "=" + strings.Join(strings.Fields(value), " ")

		}

	}

	return strings.Join(strings.Fields(raw), " ")

}

func toolCallContent(tool_call schemas.ToolCall) string {

	arguments, err := tool_call.GetArguments()

	if err == nil {

		bytes, err0 := json.MarshalIndent(arguments, "", "  ")

		if err0 == nil {
			return "```json\n" + string(bytes) + "\n```"
		}

	}

	raw := strings.TrimSpace(string(tool_call.Function.ArgumentsRaw))

	if raw == "" {
		return ""
	}

	return "```json\n" + raw + "\n```"

}

func indentTranscriptLines(lines []transcriptLine, prefix string, style utils_tty.Style) []transcriptLine {

	result := make([]transcriptLine, 0, len(lines))

	for _, line := range lines {

		if len(line.Spans) == 0 && line.HasFill == false {
			result = append(result, line)
			continue
		}

		prefix_style := style

		if line.HasFill == true {
			prefix_style = line.Fill
		}

		spans := make([]transcriptSpan, 0, len(line.Spans)+1)
		spans = append(spans, transcriptSpan{Text: prefix, Style: prefix_style})
		spans = append(spans, line.Spans...)
		line.Spans = compactSpans(spans)
		result = append(result, line)

	}

	return result

}

func muteTranscriptLines(lines []transcriptLine) []transcriptLine {

	for index := range lines {

		if lines[index].HasFill == true {
			lines[index].Fill = lines[index].Fill.WithForeground(colors.TextMuted)
		}

		for span := range lines[index].Spans {
			lines[index].Spans[span].Style = lines[index].Spans[span].Style.WithForeground(colors.TextMuted)
		}

	}

	return lines

}

func (transcript *Transcript) HandleMouse(x int, y int, rect utils_tty.Rect) bool {

	if rect.Width <= 0 || rect.Height <= 0 {
		return false
	}

	if x < rect.X || x >= rect.X+rect.Width || y < rect.Y || y >= rect.Y+rect.Height {
		return false
	}

	index := transcript.lastStart + (y - rect.Y)

	if index < 0 || index >= len(transcript.cached) {
		return false
	}

	line := transcript.cached[index]

	if line.Toggle == nil {
		return false
	}

	if transcript.expanded == nil {
		transcript.expanded = make(map[transcriptToggle]bool)
	}

	transcript.expanded[*line.Toggle] = !transcript.expanded[*line.Toggle]
	transcript.anchorIndex = index
	transcript.anchorOffset = index - transcript.lastStart
	transcript.hasAnchor = true
	transcript.cacheDirty = true

	return true

}
