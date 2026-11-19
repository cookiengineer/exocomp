package components

import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "strings"
import "time"

type footerBinding struct {
	Key   string
	Label string
}

func footerBindings(status string) []footerBinding {

	if status != "" {
		return []footerBinding{
			{Key: "Ctrl+↑↓", Label: "scroll"},
			{Key: "Ctrl+Q", Label: "quit"},
		}
	}

	return []footerBinding{
		{Key: "Ctrl+↑↓", Label: "scroll"},
		{Key: "Ctrl+S", Label: "send"},
		{Key: "/", Label: "tools"},
		{Key: "Ctrl+Q", Label: "quit"},
	}

}

func footerBindingsWidth(bindings []footerBinding) int {

	width := 0

	for index, binding := range bindings {

		if index > 0 {
			width += utils_tty.StringWidth(" · ")
		}

		width += utils_tty.StringWidth(binding.Key)
		width += 1
		width += utils_tty.StringWidth(binding.Label)

	}

	return width

}

func writeRightBindings(buffer *utils_tty.Buffer, right int, y int, minimum int, bindings []footerBinding) {

	width := footerBindingsWidth(bindings)
	available := max(0, right-minimum)
	x := max(minimum, right-width)

	if width > available {
		x = minimum
	}

	for index, binding := range bindings {

		if index > 0 {
			x += writeFooterSpan(buffer, x, y, " · ", styles.FooterHint, right)
		}

		x += writeFooterSpan(buffer, x, y, binding.Key, styles.FooterKey, right)
		x += writeFooterSpan(buffer, x, y, " ", styles.FooterHint, right)
		x += writeFooterSpan(buffer, x, y, binding.Label, styles.FooterHint, right)

		if x >= right {
			break
		}

	}

}

func writeFooterSpan(buffer *utils_tty.Buffer, x int, y int, text string, style utils_tty.Style, right int) int {

	remaining := right - x

	if remaining <= 0 {
		return 0
	}

	text = utils_tty.Truncate(text, remaining)

	if text == "" {
		return 0
	}

	buffer.WriteString(x, y, text, style)

	return utils_tty.StringWidth(text)

}

func spinnerFrame(elapsed time.Duration) string {

	if len(utils_tty.SpinnerFrames) == 0 {
		return " "
	}

	index := int(elapsed/(80*time.Millisecond)) % len(utils_tty.SpinnerFrames)

	if index < 0 {
		index = 0
	}

	return utils_tty.SpinnerFrames[index]

}

func footerActivity(status string) (string, bool) {

	trimmed := sanitizeFooterStatus(status)

	if trimmed == "" {
		return "Waiting for your instructions", false
	}

	lower := strings.ToLower(trimmed)

	if strings.Contains(lower, "waiting for the llm") || strings.Contains(lower, "session is busy") {
		return "Agent is working…", true
	}

	if strings.Contains(lower, "answering human question") {
		return "Sending your answer…", true
	}

	if strings.HasPrefix(lower, "running ") {
		return trimmed, true
	}

	return trimmed, false

}

func sanitizeFooterStatus(value string) string {

	builder := strings.Builder{}

	for _, character := range value {

		if character == '\n' || character == '\r' || character == '\t' {
			builder.WriteRune(' ')
		} else if character >= 32 && character != 127 {
			builder.WriteRune(character)
		}

	}

	return strings.Join(strings.Fields(builder.String()), " ")

}
