package fmt

import "regexp"
import "strings"

var format_code_fence_language = regexp.MustCompile("(?s)```(\\w+)?\\n(.*?)```")
var format_code_fence = regexp.MustCompile("(?s)```(.*?)```")
var format_inline_code = regexp.MustCompile("`([^`]+?)`")
var format_strong_asterisk = regexp.MustCompile(`\*\*(.*?)\*\*`)
var format_strong_underscore = regexp.MustCompile(`__(.*?)__`)
var format_emphasis_asterisk = regexp.MustCompile(`\*(.*?)\*`)
var format_emphasis_underscore = regexp.MustCompile(`_(.*?)_`)
var format_underline = regexp.MustCompile(`\+\+(.*?)\+\+`)
var format_strike = regexp.MustCompile(`~~(.*?)~~`)
var format_heading_6 = regexp.MustCompile(`(?m)^###### (.*)$`)
var format_heading_5 = regexp.MustCompile(`(?m)^##### (.*)$`)
var format_heading_4 = regexp.MustCompile(`(?m)^#### (.*)$`)
var format_heading_3 = regexp.MustCompile(`(?m)^### (.*)$`)
var format_heading_2 = regexp.MustCompile(`(?m)^## (.*)$`)
var format_heading_1 = regexp.MustCompile(`(?m)^# (.*)$`)
var format_link = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
var format_unordered_list = regexp.MustCompile(`(?m)^\s*[-*] (.*)$`)
var format_ordered_list = regexp.MustCompile(`(?m)^\s*\d+\. (.*)$`)
var format_list_group = regexp.MustCompile(`(?s)(<li>.*</li>)`)

// FormatContent converts a small markdown subset to HTML, mirroring the old
// frontend `utils/fmt/FormatContent.mjs`.
func FormatContent(raw string) string {

	tmp := raw

	tmp = strings.ReplaceAll(tmp, "&", "&amp;")
	tmp = strings.ReplaceAll(tmp, "<", "&lt;")
	tmp = strings.ReplaceAll(tmp, ">", "&gt;")

	tmp = format_code_fence_language.ReplaceAllStringFunc(tmp, func(match string) string {

		groups := format_code_fence_language.FindStringSubmatch(match)

		if len(groups) == 3 {

			lang := strings.TrimSpace(groups[1])
			code := strings.TrimSpace(groups[2])

			if lang != "" {
				return "<pre class=\"" + lang + "\">" + code + "</pre>"
			}

			return "<pre>" + code + "</pre>"

		}

		return match

	})

	tmp = format_code_fence.ReplaceAllString(tmp, "<pre><code>$1</code></pre>")
	tmp = format_inline_code.ReplaceAllString(tmp, "<code>$1</code>")

	tmp = format_strong_asterisk.ReplaceAllString(tmp, "<strong>$1</strong>")
	tmp = format_strong_underscore.ReplaceAllString(tmp, "<strong>$1</strong>")

	tmp = format_emphasis_asterisk.ReplaceAllString(tmp, "<em>$1</em>")
	tmp = format_emphasis_underscore.ReplaceAllString(tmp, "<em>$1</em>")

	tmp = format_underline.ReplaceAllString(tmp, "<u>$1</u>")
	tmp = format_strike.ReplaceAllString(tmp, "<del>$1</del>")

	tmp = format_heading_6.ReplaceAllString(tmp, "<h6>$1</h6>")
	tmp = format_heading_5.ReplaceAllString(tmp, "<h5>$1</h5>")
	tmp = format_heading_4.ReplaceAllString(tmp, "<h4>$1</h4>")
	tmp = format_heading_3.ReplaceAllString(tmp, "<h3>$1</h3>")
	tmp = format_heading_2.ReplaceAllString(tmp, "<h2>$1</h2>")
	tmp = format_heading_1.ReplaceAllString(tmp, "<h1>$1</h1>")

	tmp = format_link.ReplaceAllString(tmp, "<a href=\"$2\" target=\"_blank\" rel=\"noopener noreferrer\">$1</a>")

	tmp = format_unordered_list.ReplaceAllString(tmp, "<li>$1</li>")
	tmp = format_ordered_list.ReplaceAllString(tmp, "<li>$1</li>")
	tmp = format_list_group.ReplaceAllString(tmp, "<ul>$1</ul>")

	tmp = strings.ReplaceAll(tmp, "\n", "<br>")

	return tmp

}
