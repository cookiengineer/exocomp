package components

import "exocomp/types"
import utils_tty "exocomp/utils/tty"
import "strings"
import "testing"

func TestChoiceDialogCtrlCClearsChoices(t *testing.T) {

	dialog := NewChoiceDialog(&types.Question{
		Type:     "Choose",
		Question: "Pick one",
		Options:  []string{"A", "B", "C"},
		Multiple: true,
	})

	dialog.Selected = 2
	dialog.chosen[0] = true
	dialog.chosen[2] = true

	_, submit, cancel := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 'c', Ctrl: true})

	if dialog.Selected != 0 {
		t.Errorf("Expected selection to reset, got %d", dialog.Selected)
	}

	if len(dialog.chosen) != 0 {
		t.Errorf("Expected choices to be cleared, got %v", dialog.chosen)
	}

	if submit == true || cancel == true {
		t.Errorf("Expected Ctrl+C to only clear the choices")
	}

}

func TestChoiceDialogHeightMatchesOptions(t *testing.T) {

	dialog := NewChoiceDialog(&types.Question{
		Type:     "Choose",
		Question: "Pick one",
		Options:  []string{"A", "B", "C", "D"},
	})

	if dialog.Height(40) != 10 {
		t.Errorf("Expected height %d, got %d", 10, dialog.Height(40))
	}

}

func TestChoiceDialogRendersAllOptions(t *testing.T) {

	options := []string{"Alpha", "Beta", "Gamma"}
	dialog := NewChoiceDialog(&types.Question{
		Type:     "Choose",
		Question: "Pick one",
		Options:  options,
	})

	width := 40
	height := dialog.Height(width)
	buffer := utils_tty.NewBuffer(width, height)
	dialog.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: width, Height: height}, 0)

	builder := strings.Builder{}

	for y := 0; y < height; y++ {

		for x := 0; x < width; x++ {
			builder.WriteString(buffer.Get(x, y).Text)
		}

		builder.WriteString("\n")

	}

	rendered := builder.String()

	for _, option := range options {

		if strings.Contains(rendered, option) == false {
			t.Errorf("Expected option %q to be rendered in:\n%s", option, rendered)
		}

	}

}

func TestChoiceDialogCtrlSSubmits(t *testing.T) {

	dialog := NewChoiceDialog(&types.Question{
		Type:     "Choose",
		Question: "Pick one",
		Options:  []string{"A", "B", "C"},
	})

	answer, submit, cancel := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 's', Ctrl: true})

	if submit == false {
		t.Errorf("Expected Ctrl+S to submit")
	}

	if answer != "A" {
		t.Errorf("Expected answer %q, got %q", "A", answer)
	}

	if cancel == true {
		t.Errorf("Expected Ctrl+S not to cancel")
	}

}

func TestChoiceDialogEnterDoesNotSubmit(t *testing.T) {

	dialog := NewChoiceDialog(&types.Question{
		Type:     "Choose",
		Question: "Pick one",
		Options:  []string{"A", "B", "C"},
	})

	_, submit, cancel := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyEnter})

	if submit == true {
		t.Errorf("Expected Enter not to submit")
	}

	if cancel == true {
		t.Errorf("Expected Enter not to cancel")
	}

}
