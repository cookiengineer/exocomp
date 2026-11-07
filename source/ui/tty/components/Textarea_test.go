package components

import "testing"

func TestTextareaInsertAndNewline(t *testing.T) {

	textarea := NewTextarea()
	textarea.Insert("hello")
	textarea.InsertNewline()
	textarea.Insert("world")

	if textarea.GetText() != "hello\nworld" {
		t.Errorf("Expected hello newline world, got %q", textarea.GetText())
	}

	textarea.Backspace()

	if textarea.GetText() != "hello\nworl" {
		t.Errorf("Expected hello newline worl, got %q", textarea.GetText())
	}

}

func TestTextareaWrap(t *testing.T) {

	textarea := NewTextarea()
	textarea.SetText("one two three four")

	lines, _, _ := textarea.Wrap(7)

	if len(lines) < 2 {
		t.Errorf("Expected wrapped lines, got %d", len(lines))
	}

	if lines[0].Text != "one two" {
		t.Errorf("Expected one two, got %q", lines[0].Text)
	}

}
