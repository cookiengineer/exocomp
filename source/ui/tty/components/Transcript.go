package components

import "exocomp/engine"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "strings"
import "time"

type transcriptSpan struct {
	Text  string
	Style utils_tty.Style
}

type transcriptLine struct {
	Spans   []transcriptSpan
	HasFill bool
	Fill    utils_tty.Style
}

type Transcript struct {
	Session    *engine.Session
	Scroll     int
	start      int
	cached     []transcriptLine
	cacheWidth int
	cacheCount int
	cacheDirty bool
}

func NewTranscript(session *engine.Session) *Transcript {

	return &Transcript{
		Session:    session,
		Scroll:     0,
		cacheDirty: true,
	}

}

func (transcript *Transcript) Width(width int) int {
	return width
}

func (transcript *Transcript) Height(width int) int {
	return 0
}

func (transcript *Transcript) Invalidate() {
	transcript.cacheDirty = true
}

func (transcript *Transcript) ScrollBy(delta int) {

	transcript.Scroll += delta

	if transcript.Scroll < 0 {
		transcript.Scroll = 0
	}

}

func (transcript *Transcript) ScrollToTop() {
	transcript.Scroll = 1 << 30
}

func (transcript *Transcript) ScrollToBottom() {
	transcript.Scroll = 0
}

func (transcript *Transcript) Clear() {

	if transcript.Session != nil {
		transcript.start = len(transcript.Session.GetMessages(0))
	}

	transcript.cached = make([]transcriptLine, 0)
	transcript.Scroll = 0
	transcript.cacheWidth = 0
	transcript.cacheCount = 0
	transcript.cacheDirty = false

}

func (transcript *Transcript) Reset() {

	transcript.start = 0
	transcript.Scroll = 0
	transcript.cacheDirty = true

}

func (transcript *Transcript) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	buffer.Fill(rect, ' ', styles.Background)

	if transcript.Session == nil {
		return
	}

	messages := transcript.Session.GetMessages(0)

	if transcript.start > len(messages) {
		transcript.start = 0
	}

	messages = messages[transcript.start:]

	lines := transcript.cached

	if transcript.cacheDirty == true || transcript.cacheWidth != rect.Width || transcript.cacheCount != len(messages) {

		result := make([]transcriptLine, 0)

		for _, message := range messages {

			if message == nil {
				continue
			}

			if len(result) > 0 {
				result = append(result, transcriptLine{})
			}

			role_style := styles.TranscriptTool

			switch message.Role {

			case "user":
				role_style = styles.TranscriptUser
			case "assistant":
				role_style = styles.TranscriptAssistant
			case "tool":
				role_style = styles.TranscriptTool
			case "system":
				role_style = styles.TranscriptSystem

			}

			result = append(result, transcriptLine{
				Spans: []transcriptSpan{
					{Text: fmt.Sprintf("[%s]", message.Role), Style: role_style},
				},
			})

			content := message.Content

			if strings.TrimSpace(content) == "" && len(message.ToolCalls) > 0 {

				names := make([]string, 0)

				for _, tool_call := range message.ToolCalls {
					names = append(names, tool_call.Function.Name)
				}

				content = "Tool calls: " + strings.Join(names, ", ")

			}

			if strings.TrimSpace(content) != "" {

				rendered := renderMarkdown(content, rect.Width)
				result = append(result, rendered...)

			}

		}

		transcript.cached = result
		transcript.cacheWidth = rect.Width
		transcript.cacheCount = len(messages)
		transcript.cacheDirty = false
		lines = result

	}

	if len(lines) == 0 {

		message := "No messages yet. Ask the agent something, or type / to run a tool call."
		buffer.WriteString(rect.X+1, rect.Y, utils_tty.Truncate(message, rect.Width-2), styles.TextMuted)
		return

	}

	if transcript.Scroll > len(lines)-rect.Height {
		transcript.Scroll = max(0, len(lines)-rect.Height)
	}

	end := len(lines) - transcript.Scroll
	start := max(0, end-rect.Height)

	for row := start; row < end && row-start < rect.Height; row++ {

		line := lines[row]
		line_row := row - start

		if line.HasFill == true {
			buffer.Fill(utils_tty.Rect{X: rect.X, Y: rect.Y + line_row, Width: rect.Width, Height: 1}, ' ', line.Fill)
		}

		x := rect.X
		limit := rect.X + rect.Width

		for _, span := range line.Spans {

			if span.Text == "" {
				continue
			}

			remaining := limit - x

			if remaining <= 0 {
				break
			}

			text := utils_tty.Truncate(span.Text, remaining)

			if text == "" {
				break
			}

			buffer.WriteString(x, rect.Y+line_row, text, span.Style)
			x += utils_tty.StringWidth(text)

		}

	}

}
