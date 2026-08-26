package html

import "strings"
import "testing"

func TestParse(t *testing.T) {

	raw := `<!doctype html>
<html>
<head>
	<title>Example Page</title>
	<style>body { color: red; }</style>
	<script>alert("hello");</script>
</head>
<body>
	<h1>Welcome</h1>
	<p>Hello <strong>world</strong>, read our <a href="/about">about page</a>.</p>
	<ul>
		<li>First</li>
		<li>Second</li>
	</ul>
</body>
</html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if document == nil {
		t.Fatalf("Expected a non-nil document")
	}

	if document.Title != "Example Page" {
		t.Errorf("Expected title %q, got %q", "Example Page", document.Title)
	}

	result := document.String()

	if !strings.Contains(result, "# Welcome") {
		t.Errorf("Expected markdown to contain the h1 heading, got:\n%s", result)
	}

	if !strings.Contains(result, "Hello **world**") {
		t.Errorf("Expected markdown to contain strong text, got:\n%s", result)
	}

	if !strings.Contains(result, "[about page](https://example.com/about)") {
		t.Errorf("Expected markdown to contain a resolved link, got:\n%s", result)
	}

	if !strings.Contains(result, "- First") || !strings.Contains(result, "- Second") {
		t.Errorf("Expected markdown to contain the list, got:\n%s", result)
	}

}

func TestParse_SkipsNoise(t *testing.T) {

	raw := `<html><body><script>var x = 1;</script><style>.x {}</style><p>Visible</p></body></html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	result := document.String()

	if strings.Contains(result, "var x") {
		t.Errorf("Expected script content to be skipped, got:\n%s", result)
	}

	if strings.Contains(result, ".x {}") {
		t.Errorf("Expected style content to be skipped, got:\n%s", result)
	}

	if !strings.Contains(result, "Visible") {
		t.Errorf("Expected visible content to be kept, got:\n%s", result)
	}

}

func TestParse_CodeBlock(t *testing.T) {

	raw := `<html><body><pre><code>fmt.Println("hi")</code></pre></body></html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	result := document.String()

	if !strings.Contains(result, "```") {
		t.Errorf("Expected a code block, got:\n%s", result)
	}

	if !strings.Contains(result, `fmt.Println("hi")`) {
		t.Errorf("Expected code content, got:\n%s", result)
	}

}

func TestParse_Table(t *testing.T) {

	raw := `<html><body><table><thead><tr><th>Name</th><th>Age</th></tr></thead><tbody><tr><td>Alice</td><td>30</td></tr></tbody></table></body></html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	result := document.String()

	if !strings.Contains(result, "| Name | Age |") {
		t.Errorf("Expected a markdown table, got:\n%s", result)
	}

	if !strings.Contains(result, "| Alice | 30 |") {
		t.Errorf("Expected a markdown table row, got:\n%s", result)
	}

}

func TestParse_Text(t *testing.T) {

	raw := `<html><body><h1>Heading</h1><p>Some text</p></body></html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	text := document.Text()

	if !strings.Contains(text, "Heading") {
		t.Errorf("Expected plain text to contain heading, got:\n%s", text)
	}

	if !strings.Contains(text, "Some text") {
		t.Errorf("Expected plain text to contain body, got:\n%s", text)
	}

	if strings.Contains(text, "#") {
		t.Errorf("Expected plain text to drop markdown syntax, got:\n%s", text)
	}

}

func TestParse_ReaderMode(t *testing.T) {

	raw := `<!doctype html>
<html>
<head><title>Article</title></head>
<body>
	<nav class="nav">Home About Contact</nav>
	<div id="content" class="article-body">
		<h1>The Real Article</h1>
		<p>This is the first paragraph of actual content with enough words to matter, and a comma, and another one.</p>
		<p>This is the second paragraph, also part of the main story.</p>
	</div>
	<aside class="sidebar">Related links and advertisements</aside>
	<footer class="footer">Copyright 2026</footer>
</body>
</html>`

	document, err := Parse("https://example.com/", []byte(raw))

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	result := document.String()

	if !strings.Contains(result, "The Real Article") {
		t.Errorf("Expected the article heading, got:\n%s", result)
	}

	if !strings.Contains(result, "first paragraph") {
		t.Errorf("Expected the article body, got:\n%s", result)
	}

	if strings.Contains(result, "Home About Contact") {
		t.Errorf("Expected navigation to be stripped, got:\n%s", result)
	}

	if strings.Contains(result, "Related links") {
		t.Errorf("Expected sidebar to be stripped, got:\n%s", result)
	}

	if strings.Contains(result, "Copyright") {
		t.Errorf("Expected footer to be stripped, got:\n%s", result)
	}

}
