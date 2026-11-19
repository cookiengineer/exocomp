package components

import "exocomp/schemas"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "time"

type transcriptSpan struct {
	Text  string
	Style utils_tty.Style
}

type transcriptToggle struct {
	Message *schemas.Message
	Part    string
}

type transcriptLine struct {
	Spans     []transcriptSpan
	HasFill   bool
	Fill      utils_tty.Style
	Toggle    *transcriptToggle
	Spinner   bool
	HasBorder bool
	Border    utils_tty.Style
}

type Transcript struct {
	Scroll       int
	start        int
	messages     []*schemas.Message
	waiting      bool
	cached       []transcriptLine
	cacheWidth   int
	cacheCount   int
	cacheWaiting bool
	cacheDirty   bool
	expanded     map[transcriptToggle]bool
	lastStart    int
	lastRect     utils_tty.Rect
	anchorIndex  int
	anchorOffset int
	hasAnchor    bool
}

func NewTranscript() *Transcript {

	return &Transcript{
		Scroll:     0,
		messages:   make([]*schemas.Message, 0),
		cacheDirty: true,
		expanded:   make(map[transcriptToggle]bool),
	}

}

func (transcript *Transcript) SetMessages(messages []*schemas.Message, waiting bool) {

	transcript.messages = messages
	transcript.waiting = waiting

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

	transcript.start = len(transcript.messages)
	transcript.cached = make([]transcriptLine, 0)
	transcript.Scroll = 0
	transcript.cacheWidth = 0
	transcript.cacheCount = 0
	transcript.cacheDirty = false
	transcript.expanded = make(map[transcriptToggle]bool)
	transcript.hasAnchor = false

}

func (transcript *Transcript) Reset() {

	transcript.start = 0
	transcript.Scroll = 0
	transcript.cacheDirty = true
	transcript.expanded = make(map[transcriptToggle]bool)
	transcript.hasAnchor = false

}

func (transcript *Transcript) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	buffer.Fill(rect, ' ', styles.Background)

	messages := transcript.messages

	if transcript.start > len(messages) {
		transcript.start = 0
	}

	messages = messages[transcript.start:]

	lines := transcript.cached
	waiting := transcript.waiting
	content_width := max(1, rect.Width-2)

	if transcript.cacheDirty == true || transcript.cacheWidth != rect.Width || transcript.cacheCount != len(messages) || transcript.cacheWaiting != waiting {

		result := make([]transcriptLine, 0)

		for _, message := range messages {

			if message == nil {
				continue
			}

			if len(result) > 0 {
				result = append(result, transcriptLine{})
			}

			result = transcript.appendMessage(result, message, content_width)

		}

		if waiting == true {

			if len(result) > 0 {
				result = append(result, transcriptLine{})
			}

			result = append(result, transcriptLine{
				Spans: []transcriptSpan{
					{Text: "Thinking", Style: styles.ReasoningHeader},
				},
				Spinner:   true,
				HasBorder: true,
				Border:    styles.BorderThinking,
			})

		}

		for len(result) > 0 {

			last := result[len(result)-1]

			if len(last.Spans) == 0 && last.HasFill == false && last.Spinner == false {
				result = result[:len(result)-1]
			} else {
				break
			}

		}

		if len(result) > 0 {
			result = append(result, transcriptLine{})
		}

		transcript.cached = result
		transcript.cacheWidth = rect.Width
		transcript.cacheCount = len(messages)
		transcript.cacheWaiting = waiting
		transcript.cacheDirty = false
		lines = result

	}

	if len(lines) == 0 {

		message := "No messages yet. Ask the agent something, or type / to run a tool call."
		buffer.WriteString(rect.X+1, rect.Y, utils_tty.Truncate(message, rect.Width-2), styles.TextMuted)
		transcript.lastStart = 0
		transcript.lastRect = rect
		return

	}

	if transcript.hasAnchor == true {

		desired_start := max(0, transcript.anchorIndex-transcript.anchorOffset)
		desired_end := min(len(lines), desired_start+rect.Height)
		transcript.Scroll = max(0, len(lines)-desired_end)
		transcript.hasAnchor = false

	}

	if transcript.Scroll > len(lines)-rect.Height {
		transcript.Scroll = max(0, len(lines)-rect.Height)
	}

	end := len(lines) - transcript.Scroll
	start := max(0, end-rect.Height)
	transcript.lastStart = start
	transcript.lastRect = rect

	for row := start; row < end && row-start < rect.Height; row++ {

		line := lines[row]
		line_row := row - start
		content_x := rect.X
		line_width := rect.Width

		if line.HasBorder == true {

			buffer.Set(rect.X, rect.Y+line_row, "┃", line.Border)
			content_x = rect.X + 2
			line_width = max(1, rect.Width-2)

		}

		if line.Spinner == true {
			buffer.WriteString(content_x, rect.Y+line_row, spinnerFrame(elapsed)+" Thinking", styles.ReasoningHeader)
			continue
		}

		if line.HasFill == true {
			buffer.Fill(utils_tty.Rect{X: content_x, Y: rect.Y + line_row, Width: line_width, Height: 1}, ' ', line.Fill)
		}

		x := content_x
		limit := content_x + line_width

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
