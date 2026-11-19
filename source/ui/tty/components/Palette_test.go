package components

import "exocomp/ui/tty/styles"
import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "testing"

func newPaletteTestItems(count int) []commands.Suggestion {

	items := make([]commands.Suggestion, 0)

	for index := 0; index < count; index++ {
		items = append(items, commands.Suggestion{
			Kind:        "tool",
			Value:       fmt.Sprintf("tool-%d", index),
			Description: "test suggestion",
		})
	}

	return items

}

func TestPaletteScrollsToSelectedItem(t *testing.T) {

	palette := NewPalette()
	palette.Set(newPaletteTestItems(20))
	buffer := utils_tty.NewBuffer(40, 5)
	rect := utils_tty.Rect{X: 0, Y: 0, Width: 40, Height: 5}

	palette.Render(buffer, rect, 0)

	if palette.scroll != 0 {
		t.Errorf("Expected the initial scroll to be zero, got %d", palette.scroll)
	}

	for index := 0; index < 8; index++ {
		palette.Move(1)
	}

	palette.Render(buffer, rect, 0)

	if palette.Selected != 8 {
		t.Errorf("Expected selection 8, got %d", palette.Selected)
	}

	if palette.scroll != 4 {
		t.Errorf("Expected scroll 4, got %d", palette.scroll)
	}

	if buffer.Get(1, 4).Style != styles.Selection {
		t.Errorf("Expected the selected item to be visible")
	}

	for index := 0; index < 8; index++ {
		palette.Move(-1)
	}

	palette.Render(buffer, rect, 0)

	if palette.scroll != 0 {
		t.Errorf("Expected scroll to return to zero, got %d", palette.scroll)
	}

}

func TestPaletteHideResetsScroll(t *testing.T) {

	palette := NewPalette()
	palette.Set(newPaletteTestItems(20))
	palette.Selected = 12
	palette.scroll = 8
	palette.Hide()

	if palette.Selected != 0 {
		t.Errorf("Expected selection reset")
	}

	if palette.scroll != 0 {
		t.Errorf("Expected scroll reset")
	}

}

func TestPaletteSetClampsScroll(t *testing.T) {

	palette := NewPalette()
	palette.Set(newPaletteTestItems(20))
	palette.Selected = 12
	palette.scroll = 8
	palette.Set(newPaletteTestItems(3))

	if palette.Selected >= 3 {
		t.Errorf("Expected selection to be clamped")
	}

	if palette.scroll > 2 {
		t.Errorf("Expected scroll to be clamped")
	}

}
