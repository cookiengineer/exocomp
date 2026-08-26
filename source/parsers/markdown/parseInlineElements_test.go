package markdown

import "testing"

func TestParseInlineElements(t *testing.T) {

	t.Run("parseInlineElements(emoji)", func(t *testing.T) {

		children1, err1 := parseInlineElements("before :balloon: after", 1, false)

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if len(children1) == 1 {

			if children1[0].Type != "#text" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children1[0].Type, "#text")
			}

			if children1[0].Text != "before 🎈 after" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children1[0].Text, "before 🎈 after")
			}

		} else {
			t.Errorf("Expected %d elements to be %d", len(children1), 1)
		}

		children2, err2 := parseInlineElements("before :balloon after", 1, false)

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if len(children2) == 1 {

			if children2[0].Type != "#text" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children2[0].Type, "#text")
			}

			if children2[0].Text != "before :balloon after" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children2[0].Text, "before :balloon after")
			}

		} else {
			t.Errorf("Expected %d elements to be %d", len(children2), 1)
		}

	})

	t.Run("parseInlineElements(code)", func(t *testing.T) {

		children, err := parseInlineElements("before `foo(\"bar\")` after", 1, false)

		if err != nil {
			t.Errorf("Expected %v to be nil", err)
		}

		if len(children) == 3 {

			if children[0].Type != "#text" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[0].Type, "#text")
			}

			if children[0].Text != "before" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[0].Text, "before")
			}

			if children[1].Type != "code" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[1].Type, "code")
			}

			if children[1].Text != "foo(\"bar\")" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[1].Text, "foo(\"bar\")")
			}

			if children[2].Type != "#text" {
				t.Errorf("Expected %s to be %s", children[2].Type, "#text")
			}

			if children[2].Text != "after" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[2].Text, "after")
			}

		} else {
			t.Errorf("Expected %d elements to be %d", len(children), 3)
		}

	})

	t.Run("parseInlineElements(code_with_colon)", func(t *testing.T) {

		children, err := parseInlineElements("How it works: `program parameter /foo:bar`", 1, false)

		if err != nil {
			t.Errorf("Expected %v to be nil", err)
		}

		if len(children) == 2 {

			if children[0].Type != "#text" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[0].Type, "#text")
			}

			if children[0].Text != "How it works:" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[0].Text, "How it works:")
			}

			if children[1].Type != "code" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[1].Type, "code")
			}

			if children[1].Text != "program parameter /foo:bar" {
				t.Errorf("Expected \"%s\" to be \"%s\"", children[1].Text, "program parameter /foo:bar")
			}

		} else {
			t.Errorf("Expected %d elements to be %d", len(children), 2)
		}

	})

}
