package fmt

import "strings"
import "testing"

func TestFormatContent_EscapesHTML(t *testing.T) {

	result := FormatContent("<script>alert(1)</script>")

	if strings.Contains(result, "<script>") == true {
		t.Errorf("Expected HTML to be escaped, got %q", result)
	}

}

func TestFormatContent_Bold(t *testing.T) {

	result := FormatContent("**bold**")

	if result != "<strong>bold</strong>" {
		t.Errorf("Expected %q, got %q", "<strong>bold</strong>", result)
	}

}

func TestFormatContent_InlineCode(t *testing.T) {

	result := FormatContent("`code`")

	if result != "<code>code</code>" {
		t.Errorf("Expected %q, got %q", "<code>code</code>", result)
	}

}

func TestSanitizeContent(t *testing.T) {

	result := SanitizeContent("<b>\"quoted\"</b>")

	if result != "&lt;b&gt;&quot;quoted&quot;&lt;/b&gt;" {
		t.Errorf("Expected escaped entities, got %q", result)
	}

}
