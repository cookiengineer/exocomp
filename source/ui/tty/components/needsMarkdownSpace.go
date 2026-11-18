package components

import "strings"

func needsMarkdownSpace(previous []transcriptSpan, next []transcriptSpan) bool {

	if len(previous) == 0 || len(next) == 0 {
		return false
	}

	last_span := previous[len(previous)-1]
	first_span := next[0]

	if last_span.Text == "" || first_span.Text == "" {
		return false
	}

	last := last_span.Text[len(last_span.Text)-1]
	first := first_span.Text[0]

	if last == ' ' || last == '\n' || last == '\t' {
		return false
	}

	if first == ' ' || first == '\n' || first == '\t' {
		return false
	}

	return strings.ContainsRune(".,;:!?)]}", rune(first)) == false

}
