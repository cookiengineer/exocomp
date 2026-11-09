package components

func compactSpans(spans []transcriptSpan) []transcriptSpan {

	result := make([]transcriptSpan, 0)

	for _, span := range spans {

		if span.Text == "" {
			continue
		}

		if len(result) > 0 {

			last := result[len(result)-1]

			if last.Style.Equals(span.Style) == true {
				result[len(result)-1].Text += span.Text
				continue
			}

		}

		result = append(result, span)

	}

	return result

}
