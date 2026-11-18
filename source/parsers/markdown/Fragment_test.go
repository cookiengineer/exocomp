package markdown

import "testing"

func TestParseFragment(t *testing.T) {

	document := ParseFragment("# Heading\n\n- first item\n- second item\n\n```go\nfmt.Println(\"x\")\n```\n\nThis is **bold** and `code`.")

	if document == nil {
		t.Fatal("Expected a document")
	}

	if len(document.Errors) != 0 {
		t.Fatalf("Expected no errors, got %v", document.Errors)
	}

	if len(document.Body) == 0 {
		t.Fatal("Expected a document body")
	}

	if document.Body[0].Type != "h1" {
		t.Errorf("Expected h1, got %q", document.Body[0].Type)
	}

}

func TestParseFragmentExtendedSyntax(t *testing.T) {

	document := ParseFragment("###### Heading\n\n> quoted\n\n---\n\n~~deleted~~ and ++underlined++")

	if document == nil {
		t.Fatal("Expected a document")
	}

	if len(document.Errors) != 0 {
		t.Fatalf("Expected no errors, got %v", document.Errors)
	}

	found_h6 := false
	found_blockquote := false
	found_hr := false
	found_del := false
	found_underline := false

	for _, element := range document.Body {

		switch element.Type {

		case "h6":
			found_h6 = true
		case "blockquote":
			found_blockquote = true
		case "hr":
			found_hr = true
		case "p":
			for _, child := range element.Children {
				if child.Type == "del" {
					found_del = true
				}
				if child.Type == "u" {
					found_underline = true
				}
			}
		}

	}

	if found_h6 == false {
		t.Errorf("Expected h6 element")
	}

	if found_blockquote == false {
		t.Errorf("Expected blockquote element")
	}

	if found_hr == false {
		t.Errorf("Expected horizontal rule element")
	}

	if found_del == false {
		t.Errorf("Expected deleted element")
	}

	if found_underline == false {
		t.Errorf("Expected underline element")
	}

}

func TestParseFragmentEmptyCodeFence(t *testing.T) {

	document := ParseFragment("```\ncode\n```")

	if document == nil {
		t.Fatal("Expected a document")
	}

	if len(document.Errors) != 0 {
		t.Fatalf("Expected no errors, got %v", document.Errors)
	}

	if len(document.Body) == 0 || document.Body[0].Type != "pre" {
		t.Fatalf("Expected pre element, got %v", document.Body)
	}

	if document.Body[0].GetAttribute("class") != "" {
		t.Errorf("Expected empty class, got %q", document.Body[0].GetAttribute("class"))
	}

}

func TestParseFragmentRawHTMLIsText(t *testing.T) {

	document := ParseFragment("<script>alert(1)</script>")

	if document == nil {
		t.Fatal("Expected a document")
	}

	if len(document.Body) != 1 || document.Body[0].Type != "p" {
		t.Fatalf("Expected a paragraph, got %v", document.Body)
	}

	if len(document.Body[0].Children) != 1 || document.Body[0].Children[0].Type != "#text" {
		t.Fatalf("Expected raw HTML text, got %v", document.Body[0].Children)
	}

	if document.Body[0].Children[0].Text != "<script>alert(1)</script>" {
		t.Errorf("Expected raw HTML text, got %q", document.Body[0].Children[0].Text)
	}

}

func TestParseFragmentLinksAndImages(t *testing.T) {

	document := ParseFragment("[example](https://example.com/) and ![alt](image.png)")

	if document == nil {
		t.Fatal("Expected a document")
	}

	if len(document.Errors) != 0 {
		t.Fatalf("Expected no errors, got %v", document.Errors)
	}

	if len(document.Body) != 1 || document.Body[0].Type != "p" {
		t.Fatalf("Expected a paragraph, got %v", document.Body)
	}

	found_link := false
	found_image := false

	for _, child := range document.Body[0].Children {

		if child.Type == "a" && child.GetAttribute("href") == "https://example.com/" {
			found_link = true
		}

		if child.Type == "img" && child.GetAttribute("src") == "image.png" {
			found_image = true
		}

	}

	if found_link == false {
		t.Errorf("Expected link element")
	}

	if found_image == false {
		t.Errorf("Expected image element")
	}

}
