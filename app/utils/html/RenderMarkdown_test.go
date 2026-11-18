package html

import "strings"
import "testing"

func TestRenderMarkdownEscapesHTML(t *testing.T) {

	result := RenderMarkdown("<script>alert(1)</script>")

	if strings.Contains(result, "<script>") == true {
		t.Errorf("Expected HTML to be escaped, got %q", result)
	}

	if strings.Contains(result, "&lt;script&gt;") == false {
		t.Errorf("Expected escaped script text, got %q", result)
	}

}

func TestRenderMarkdownInlineFormatting(t *testing.T) {

	result := RenderMarkdown("This is **bold**, *emphasized*, ++underlined++, ~~deleted~~ and `code`.")

	if strings.Contains(result, "<strong>bold</strong>") == false {
		t.Errorf("Expected bold element, got %q", result)
	}

	if strings.Contains(result, "<em>emphasized</em>") == false {
		t.Errorf("Expected emphasis element, got %q", result)
	}

	if strings.Contains(result, "<u>underlined</u>") == false {
		t.Errorf("Expected underline element, got %q", result)
	}

	if strings.Contains(result, "<del>deleted</del>") == false {
		t.Errorf("Expected deleted element, got %q", result)
	}

	if strings.Contains(result, "<code>code</code>") == false {
		t.Errorf("Expected code element, got %q", result)
	}

}

func TestRenderMarkdownBlocks(t *testing.T) {

	result := RenderMarkdown("# Heading\n\n- first item\n- second item\n\n```go\nfmt.Println(\"x\")\n```\n\n> quoted\n\n---")

	if strings.Contains(result, "<h1>Heading</h1>") == false {
		t.Errorf("Expected heading, got %q", result)
	}

	if strings.Contains(result, "<ul><li>first item</li><li>second item</li></ul>") == false {
		t.Errorf("Expected list, got %q", result)
	}

	if strings.Contains(result, "<pre class=\"go\">fmt.Println(&#34;x&#34;)</pre>") == false {
		t.Errorf("Expected code block, got %q", result)
	}

	if strings.Contains(result, "<blockquote><p>quoted</p></blockquote>") == false {
		t.Errorf("Expected blockquote, got %q", result)
	}

	if strings.Contains(result, "<hr>") == false {
		t.Errorf("Expected horizontal rule, got %q", result)
	}

}

func TestRenderMarkdownLinks(t *testing.T) {

	result := RenderMarkdown("[example](https://example.com/)")

	if strings.Contains(result, "href=\"https://example.com/\"") == false {
		t.Errorf("Expected link href, got %q", result)
	}

	if strings.Contains(result, "target=\"_blank\"") == false {
		t.Errorf("Expected target attribute, got %q", result)
	}

	unsafe := RenderMarkdown("[bad](javascript:alert(1))")

	if strings.Contains(unsafe, "javascript:") == true {
		t.Errorf("Expected unsafe URL to be removed, got %q", unsafe)
	}

}

func TestRenderMarkdownLineBreak(t *testing.T) {

	result := RenderMarkdown("First\nSecond")

	if strings.Contains(result, "First<br>Second") == false {
		t.Errorf("Expected line break, got %q", result)
	}

}
