package components

import "strings"

func (renderer *markdownRenderer) inline(text string, default_kind string) []transcriptSpan {

	result := make([]transcriptSpan, 0)
	default_style := markdownStyle(renderer.styles, default_kind)
	builder := strings.Builder{}

	flush := func() {

		if builder.Len() > 0 {
			result = append(result, transcriptSpan{Text: builder.String(), Style: default_style})
			builder.Reset()
		}

	}

	for index := 0; index < len(text); {

		if text[index] == '`' {

			end := strings.Index(text[index+1:], "`")

			if end >= 0 {
				flush()
				inner := text[index+1 : index+1+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "code")})
				index = index + end + 2
				continue
			}

		}

		if strings.HasPrefix(text[index:], "**") == true {

			end := strings.Index(text[index+2:], "**")

			if end >= 0 {
				flush()
				inner := text[index+2 : index+2+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "strong")})
				index = index + end + 4
				continue
			}

		}

		if strings.HasPrefix(text[index:], "~~") == true {

			end := strings.Index(text[index+2:], "~~")

			if end >= 0 {
				flush()
				inner := text[index+2 : index+2+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "strike")})
				index = index + end + 4
				continue
			}

		}

		if text[index] == '*' {

			end := strings.Index(text[index+1:], "*")

			if end > 0 {
				flush()
				inner := text[index+1 : index+1+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "emphasis")})
				index = index + end + 2
				continue
			}

		}

		if strings.HasPrefix(text[index:], "![") == true {

			label, url, end, ok := parseMarkdownLink(text, index, true)

			if ok == true {
				flush()
				result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "imageText")})

				if url != "" {
					result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
				}

				index = end
				continue
			}

		}

		if text[index] == '[' {

			label, url, end, ok := parseMarkdownLink(text, index, false)

			if ok == true {
				flush()
				result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "link")})

				if url != "" && url != label {
					result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
				}

				index = end
				continue
			}

		}

		builder.WriteByte(text[index])
		index++

	}

	flush()

	return result

}
