package components

import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "regexp"
import "strings"

var markdownHeading = regexp.MustCompile(`^(#{1,6})[ \t]+(.*)$`)
var markdownUnordered = regexp.MustCompile(`^([ \t]*)([-*+])[ \t]+(.*)$`)
var markdownOrdered = regexp.MustCompile(`^([ \t]*)([0-9]+)[.)][ \t]+(.*)$`)

func renderMarkdown(text string, width int) []transcriptLine {

	renderer := &markdownRenderer{
		width:  max(1, width),
		styles: make(map[string]utils_tty.Style),
		lines:  make([]transcriptLine, 0),
	}

	renderer.styles["text"] = styles.MarkdownText
	renderer.styles["heading"] = styles.MarkdownHeading
	renderer.styles["strong"] = styles.MarkdownStrong
	renderer.styles["emphasis"] = styles.MarkdownEmphasis
	renderer.styles["code"] = styles.MarkdownCode
	renderer.styles["codeBlock"] = styles.MarkdownCodeBlock
	renderer.styles["link"] = styles.MarkdownLink
	renderer.styles["linkUrl"] = styles.MarkdownLinkURL
	renderer.styles["blockQuote"] = styles.MarkdownBlockQuote
	renderer.styles["listItem"] = styles.MarkdownListItem
	renderer.styles["listEnumeration"] = styles.MarkdownListEnumeration
	renderer.styles["horizontalRule"] = styles.MarkdownHorizontalRule
	renderer.styles["image"] = styles.MarkdownImage
	renderer.styles["imageText"] = styles.MarkdownImageText
	renderer.styles["strike"] = styles.MarkdownStrike

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	lines := strings.Split(normalized, "\n")
	in_code := false

	for _, raw := range lines {

		trimmed := strings.TrimSpace(raw)

		if in_code == true {

			if strings.HasPrefix(trimmed, "```") == true {

				in_code = false
				renderer.lines = append(renderer.lines, transcriptLine{})

			} else {

				prefix := []transcriptSpan{
					{Text: "  ", Style: markdownStyle(renderer.styles, "codeBlock")},
				}
				spans := []transcriptSpan{
					{Text: raw, Style: markdownStyle(renderer.styles, "codeBlock")},
				}
				renderer.appendWrapped(prefix, prefix, spans, true)

			}

			continue

		}

		if strings.HasPrefix(trimmed, "```") == true {

			in_code = true
			continue

		}

		if trimmed == "" {

			renderer.lines = append(renderer.lines, transcriptLine{})
			continue

		}

		is_rule := false

		if len(trimmed) >= 3 {

			without_spaces := strings.ReplaceAll(trimmed, " ", "")

			if len(without_spaces) >= 3 {

				first := without_spaces[0]

				if first == '-' || first == '*' || first == '_' {

					is_rule = true

					for index := 1; index < len(without_spaces); index++ {

						if without_spaces[index] != first {

							is_rule = false
							break

						}

					}

				}

			}

		}

		if is_rule == true {

			renderer.lines = append(renderer.lines, transcriptLine{
				Spans: []transcriptSpan{
					{Text: strings.Repeat("─", renderer.width), Style: markdownStyle(renderer.styles, "horizontalRule")},
				},
			})
			continue

		}

		if heading := markdownHeading.FindStringSubmatch(trimmed); heading != nil {

			spans := renderer.inline(strings.TrimSpace(heading[2]), "heading")
			renderer.lines = append(renderer.lines, transcriptLine{Spans: compactSpans(spans)})
			continue

		}

		if strings.HasPrefix(trimmed, ">") == true {

			prefix := []transcriptSpan{
				{Text: "│ ", Style: markdownStyle(renderer.styles, "blockQuote")},
			}
			spans := renderer.inline(strings.TrimSpace(trimmed[1:]), "blockQuote")
			renderer.appendWrapped(prefix, prefix, spans, false)
			continue

		}

		if item := markdownUnordered.FindStringSubmatch(raw); item != nil {

			renderer.addList(item[1], "•", strings.TrimSpace(item[3]), false)
			continue

		}

		if item := markdownOrdered.FindStringSubmatch(raw); item != nil {

			renderer.addList(item[1], item[2]+".", strings.TrimSpace(item[3]), true)
			continue

		}

		spans := renderer.inline(strings.TrimSpace(raw), "text")
		renderer.appendWrapped(nil, nil, spans, false)

	}

	return renderer.lines

}
