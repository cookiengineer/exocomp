package markdown

import utils_strings "exocomp/utils/strings"
import "strings"
import "testing"

func TestDocument(t *testing.T) {

	t.Run("Parse(header)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- title: Example Title",
			"- summary: This is an example article description. It is not allowed to have newlines.",
			"- date: 2025-12-31",
			"- tags: software, network, privacy",
			"===",
			"",
			"### Example Title",
			"",
			"This is the first paragraph.",
			"",
		}, "\n"))

		document := NewDocument("/tmp/examples/tests/header.md")
		errors := document.Parse(bytes)

		if document.IsValid() == false {
			t.Errorf("Expected %v to be %v", document.IsValid(), true)
		}

		if len(errors) != 0 {
			t.Errorf("Expected %d errors to be %d", len(errors), 0)
		}

		if len(document.Body) == 3 {

			if document.Body[0].Type != "h3" {
				t.Errorf("Expected %s to be %s", document.Body[0].Type, "h3")
			}

			if document.Body[0].GetAttribute("id") != "example-title" {
				t.Errorf("Expected id %s to be %s", document.Body[0].GetAttribute("id"), "example-title")
			}

			if len(document.Body[0].Children) == 1 {

				if document.Body[0].Children[0].Type != "#text" {
					t.Errorf("Expected %s to be %s", document.Body[0].Children[0].Type, "#text")
				}

				if document.Body[0].Children[0].Text != "Example Title" {
					t.Errorf("Expected \"%s\" to be \"%s\"", document.Body[0].Children[0].Text, "Example Title")
				}

			} else {
				t.Errorf("Expected %d elements to be %d", len(document.Body[0].Children), 1)
			}

			if document.Body[1].Type != "#text" {
				t.Errorf("Expected %s to be %s", document.Body[1].Type, "#text")
			}

			if document.Body[2].Type != "p" {
				t.Errorf("Expected %s to be %s", document.Body[2].Type, "p")
			}

		} else {
			t.Errorf("Expected %d elements to be %d", len(document.Body), 3)
		}

	})

	t.Run("Parse(invalid header)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- title: Example Title",
			"- date: 2025-12-31",
			"- tags: software, network, privacy",
			"===",
			"",
			"### Example Title",
			"",
			"This is the first paragraph.",
			"",
		}, "\n"))

		document := NewDocument("/tmp/examples/tests/invalid-header.md")
		errors := document.Parse(bytes)

		if document.IsValid() == true {
			t.Errorf("Expected %v to be %v", document.IsValid(), false)
		}

		if len(errors) == 1 {

			expected := "Line 1: Missing \"summary\" (string) frontmatter field"

			if errors[0].Error() != expected {
				t.Errorf("Expected \"%s\" instead of \"%s\" error", expected, errors[0].Error())
			}

		} else {
			t.Errorf("Expected %d errors to be %d", len(errors), 1)
		}

	})

	t.Run("Parse(inline elements)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- date: 2025-12-31",
			"- title: Example Title",
			"- summary: This is an example article description.",
			"- tags: software, network, privacy",
			"===",
			"",
			"### Example Title",
			"",
			"This is the first paragraph with `code`, some *emphasized* and **bold** ~and deleted~ text.",
			"",
		}, "\n"))

		document := NewDocument("/tmp/examples/tests/inline-elements.md")
		errors := document.Parse(bytes)

		if document.IsValid() != true {
			t.Errorf("Expected %v to be %v", document.IsValid(), true)
		}

		if len(errors) == 0 {

			if len(document.Body) == 3 {

				if document.Body[0].Type != "h3" {
					t.Errorf("Expected %s to be %s", document.Body[0].Type, "h3")
				}

				if document.Body[0].GetAttribute("id") != "example-title" {
					t.Errorf("Expected id %s to be %s", document.Body[0].GetAttribute("id"), "example-title")
				}

				if document.Body[1].Type != "#text" {
					t.Errorf("Expected %s to be %s", document.Body[1].Type, "#text")
				}

				if document.Body[2].Type == "p" {

					if len(document.Body[2].Children) == 9 {

						if document.Body[2].Children[0].Type != "#text" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[0].Type, "#text")
						}

						if document.Body[2].Children[0].Text != "This is the first paragraph with" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[0].Text, "This is the first paragraph with")
						}

						if document.Body[2].Children[1].Type != "code" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[1].Type, "code")
						}

						if document.Body[2].Children[1].Text != "code" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[1].Text, "code")
						}

						if document.Body[2].Children[2].Type != "#text" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[1].Type, "#text")
						}

						if document.Body[2].Children[2].Text != ", some" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[2].Text, ", some")
						}

						if document.Body[2].Children[3].Type != "em" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[3].Type, "em")
						}

						if document.Body[2].Children[3].Text != "emphasized" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[3].Text, "emphasized")
						}

						if document.Body[2].Children[4].Type != "#text" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[4].Type, "#text")
						}

						if document.Body[2].Children[4].Text != "and" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[4].Text, "and")
						}

						if document.Body[2].Children[5].Type != "b" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[5].Type, "b")
						}

						if document.Body[2].Children[5].Text != "bold" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[5].Text, "bold")
						}

						if document.Body[2].Children[6].Type != "#text" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[6].Type, "#text")
						}

						if document.Body[2].Children[6].Text != "" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[6].Text, "")
						}

						if document.Body[2].Children[7].Type != "del" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[7].Type, "del")
						}

						if document.Body[2].Children[7].Text != "and deleted" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[7].Text, "and deleted")
						}

						if document.Body[2].Children[8].Type != "#text" {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[8].Type, "#text")
						}

						if document.Body[2].Children[8].Text != "text." {
							t.Errorf("Expected %s to be %s", document.Body[2].Children[8].Text, "text.")
						}

					} else {
						t.Errorf("Expected %d elements to be %d", len(document.Body[2].Children), 9)
					}

				} else {
					t.Errorf("Expected %s to be %s", document.Body[2].Type, "p")
				}

			} else {
				t.Errorf("Expected %d elements to be %d", len(document.Body), 3)
			}

		} else {
			t.Errorf("Expected %d errors to be %d", len(errors), 0)
		}

	})

	t.Run("Parse(images)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- title: Example Title",
			"- summary: This is an example article description.",
			"- date: 2025-12-31",
			"- tags: software, network, privacy",
			"===",
			"",
			"![Relative Image 1](./path/to/image.png)",
			"![Relative Image 2](../path/to/image.png)",
			"![Relative Image 3](../../path/to/image.png)",
			"![Relative Image 4](../../../path/to/image.png)",
			"",
		}, "\n"))

		document := NewDocument("/articles/examples/images.md")
		errors := document.Parse(bytes)

		if len(errors) == 0 {

			if len(document.Body) == 4 {

				image1 := document.Body[0].RenderInto(document, "")

				if image1 != "<img alt=\"Relative Image 1\" src=\"/articles/examples/path/to/image.png\"/>" {
					t.Errorf("Expected %s to resolve to %s", document.Body[0].GetAttribute("src"), "/articles/examples/path/to/image.png")
				}

				image2 := document.Body[1].RenderInto(document, "")

				if image2 != "<img alt=\"Relative Image 2\" src=\"/articles/path/to/image.png\"/>" {
					t.Errorf("Expected %s to resolve to %s", document.Body[1].GetAttribute("src"), "/articles/path/to/image.png")
				}

				image3 := document.Body[2].RenderInto(document, "")

				if image3 != "<img alt=\"Relative Image 3\" src=\"/path/to/image.png\"/>" {
					t.Errorf("Expected %s to resolve to %s", document.Body[2].GetAttribute("src"), "/path/to/image.png")
				}

				image4 := document.Body[3].RenderInto(document, "")

				if image4 != "<img alt=\"Relative Image 4\" src=\"/path/to/image.png\"/>" {
					t.Errorf("Expected %s to resolve to %s", document.Body[3].GetAttribute("src"), "/path/to/image.png")
				}

			} else {
				t.Errorf("Expected %d elements to be %d", len(document.Body), 4)
			}

		} else {
			t.Errorf("Expected %d errors to be %d", len(errors), 0)
		}

	})

	t.Run("Parse(links)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- title: Example Title",
			"- summary: This is an example article description.",
			"- date: 2025-12-31",
			"- tags: software, network, privacy",
			"===",
			"",
			"[Relative Link 1](./path/to/file.html)",
			"[Relative Link 2](../path/to/file.html)",
			"[Relative Link 3](../../path/to/file.html)",
			"[Relative Link 4](../../../path/to/file.html)",
			"[Absolute Link 1](https://example.com/index.html)",
			"[Absolute Link 2](https://example.com/path/to/../file.html)",
			"",
		}, "\n"))

		document := NewDocument("/articles/examples/links.md")
		errors := document.Parse(bytes)

		if len(errors) == 0 {

			if len(document.Body) == 1 {

				if document.Body[0].Type == "p" {

					if len(document.Body[0].Children) == 6 {

						link1 := document.Body[0].Children[0].RenderInto(document, "")

						if link1 != "<a href=\"/articles/examples/path/to/file.html\">Relative Link 1</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[0].GetAttribute("href"), "/articles/examples/path/to/file.html")
						}

						link2 := document.Body[0].Children[1].RenderInto(document, "")

						if link2 != "<a href=\"/articles/path/to/file.html\">Relative Link 2</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[1].GetAttribute("href"), "/articles/path/to/file.html")
						}

						link3 := document.Body[0].Children[2].RenderInto(document, "")

						if link3 != "<a href=\"/path/to/file.html\">Relative Link 3</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[2].GetAttribute("href"), "/path/to/file.html")
						}

						link4 := document.Body[0].Children[3].RenderInto(document, "")

						if link4 != "<a href=\"/path/to/file.html\">Relative Link 4</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[3].GetAttribute("href"), "/path/to/file.html")
						}

						link5 := document.Body[0].Children[4].RenderInto(document, "")

						if link5 != "<a class=\"icon-website\" href=\"https://example.com/index.html\" target=\"_blank\">Absolute Link 1</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[4].GetAttribute("href"), "https://example.com/index.html")
						}

						link6 := document.Body[0].Children[5].RenderInto(document, "")

						if link6 != "<a class=\"icon-website\" href=\"https://example.com/path/file.html\" target=\"_blank\">Absolute Link 2</a>" {
							t.Errorf("Expected %s to resolve to %s", document.Body[0].Children[5].GetAttribute("href"), "https://example.com/path/file.html")
						}

					} else {
						t.Errorf("Expected %d children to be %d", len(document.Body[0].Children), 6)
					}

				} else {
					t.Errorf("Expected %s to be %s", document.Body[0].Type, "p")
				}

			} else {
				t.Errorf("Expected %d elements to be %d", len(document.Body), 4)
			}

		} else {
			t.Errorf("Expected %d errors to be %d", len(errors), 0)
		}

	})

	t.Run("String(article)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"===",
			"- title: Example Article",
			"- summary: This is an example article description. It is not allowed to have newlines.",
			"- date: 2025-12-31",
			"- tags: network, privacy, software",
			"- image: /articles/example-article.jpg",
			"===",
			"",
			"## Example Article",
			"",
			"This is an example article teasers. It is allowed to have newlines, `code` and ~deleted~ elements,",
			"and **bold** and *emphasized* content.",
			"",
			"### Example Headline",
			"",
			"This is another section of the article.",
			"",
			"- This is a list",
			"- [with](http://some.example.com/links.html) and descriptions",
			"- and with `code` and ~deleted~ and **bold** lines",
			"",
			"```bash",
			"if [[ ! -d ~/.config/example ]]; then",
			"\texample(true, \"hello\", \"world!\");",
			"fi;",
			"```",
			"",
			"#### Example Sub Headline",
			"",
			"This is another paragraph of the article, but this one has a very very very very very long text",
			"content which hopefully will be automatically split at 100 characters. The paragraph grows even",
			"longer because it has multiple sentences.",
			"",
			"#### Another Sub Headline",
			"",
			"| First       |    Second     |                                            Third |",
			"|:------------|:-------------:|-------------------------------------------------:|",
			"| Thing Three |  Feature Two  |                                Cell with content |",
			"| Thing One   | Feature Three |                Another Cell with shorter content |",
			"| Thing Two   |  Feature One  | Another Cell with way ~shorter~ *longer* content |",
			"",
		}, "\n"))

		document := NewDocument("/articles/example-article.md")
		errors := document.Parse(bytes)

		if document.IsValid() == false {
			t.Errorf("Expected %v to be %v", document.IsValid(), true)
		}

		if len(errors) != 0 {
			t.Errorf("Expected %d errors to be %d", len(errors), 0)
		}

		got := document.String()
		want := string(bytes)

		if got != want {

			lines_want := strings.Split(want, "\n")
			lines_got := strings.Split(got, "\n")

			diff := utils_strings.Diff(lines_want, lines_got)

			if diff != "" {
				t.Errorf("Markdown document mismatch (-want +got):\n%s", diff)
			}

		}

	})

}
