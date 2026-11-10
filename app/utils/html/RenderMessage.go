package html

import app_fmt "exocomp-app/utils/fmt"
import "exocomp-app/schemas"
import "fmt"
import "strings"

func RenderMessage(message *schemas.Message, with_empty_content bool) string {

	if message == nil {
		return ""
	}

	if message.Role == "assistant" {

		has_content := message.Content != ""
		has_reasoning := message.ReasoningContent != ""

		if has_content == true || has_reasoning == true {

			html := ""

			if has_reasoning == true {
				html += "<details class=\"reasoning\">"
				html += "<summary>Thought</summary>"
				html += app_fmt.FormatContent(message.ReasoningContent)
				html += "</details>"
			}

			if has_content == true {
				html += app_fmt.FormatContent(message.Content)
			}

			return html

		} else if with_empty_content == true {
			return "(no content)"
		}

	} else if message.Role == "system" {

		if message.Content != "" {
			return app_fmt.FormatContent(message.Content)
		} else if with_empty_content == true {
			return "(no content)"
		}

	} else if message.Role == "tool" {

		tmp := strings.Split(message.Content, "\n")

		if len(tmp) == 1 {

			return fmt.Sprintf("<pre>%s</pre>", tmp[0])

		} else if len(tmp) > 1 {

			return strings.Join([]string{
				"<details>",
				"<summary>",
				fmt.Sprintf("<pre>%s</pre>", strings.TrimSpace(tmp[0])),
				"</summary>",
				fmt.Sprintf("<pre>%s</pre>", strings.Join(tmp[1:], "\n")),
				"</details>",
			}, "")

		}

		return "<pre>(no content)</pre>"

	} else if message.Role == "user" {

		if message.Content != "" {
			return app_fmt.FormatContent(message.Content)
		} else if with_empty_content == true {
			return "(no content)"
		}

	}

	return ""

}
