package components

import "exocomp/ui/tty/colors"
import utils_tty "exocomp/utils/tty"
import "strings"
import "testing"

func markdownText(line transcriptLine) string {

	builder := strings.Builder{}

	for _, span := range line.Spans {
		builder.WriteString(span.Text)
	}

	return builder.String()

}

func TestRenderMarkdownHeadingAndList(t *testing.T) {

	lines := renderMarkdown("# Heading\n\n- first item\n- second item", 60)

	found_heading := false
	found_list := false

	for _, line := range lines {

		text := markdownText(line)

		if text == "Heading" {

			found_heading = true

			if len(line.Spans) == 0 || line.Spans[0].Style.Attributes&utils_tty.AttrBold == 0 {
				t.Errorf("Expected bold heading style")
			}

		}

		if strings.Contains(text, "• first item") {
			found_list = true
		}

	}

	if found_heading == false {
		t.Errorf("Expected a heading line")
	}

	if found_list == false {
		t.Errorf("Expected a list item line")
	}

}

func TestRenderMarkdownCodeAndBold(t *testing.T) {

	lines := renderMarkdown("```go\nfmt.Println(\"x\")\n```\n\nThis is **bold** and `code`.", 60)

	found_code := false
	found_bold := false
	found_inline := false

	for _, line := range lines {

		text := markdownText(line)

		if strings.Contains(text, "fmt.Println") {

			found_code = true

			if line.HasFill == false {
				t.Errorf("Expected code block fill")
			}

		}

		for _, span := range line.Spans {

			if span.Text == "bold" && span.Style.Attributes&utils_tty.AttrBold != 0 {
				found_bold = true
			}

			if span.Text == "code" && span.Style.Foreground == colors.MarkdownCode {
				found_inline = true
			}

		}

	}

	if found_code == false {
		t.Errorf("Expected a code block line")
	}

	if found_bold == false {
		t.Errorf("Expected bold inline text")
	}

	if found_inline == false {
		t.Errorf("Expected inline code style")
	}

}

func TestRenderMarkdownPreservesBlankLines(t *testing.T) {

	lines := renderMarkdown("First\n\n\nSecond", 40)
	blanks := 0

	for _, line := range lines {
		if len(line.Spans) == 0 && line.HasFill == false {
			blanks++
		}
	}

	if blanks == 2 {
		// Expected two preserved blank lines.
	} else {
		t.Errorf("Expected 2 blank lines, got %d", blanks)
	}

}
