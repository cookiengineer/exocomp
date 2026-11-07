package components

import "exocomp/types"
import utils_tty "exocomp/utils/tty"
import "testing"

func TestQuestionDialogHeightSingleLine(t *testing.T) {

	dialog := NewQuestionDialog(&types.Question{Type: "Ask", Question: "What is your name?"})

	if dialog.Height(64) != 8 {
		t.Errorf("Expected height %d, got %d", 8, dialog.Height(64))
	}

}

func TestQuestionDialogHeightGrowsWithNewlines(t *testing.T) {

	dialog := NewQuestionDialog(&types.Question{Type: "Ask", Question: "What is your name?"})
	dialog.Input.SetText("a\nb\nc\nd")

	if dialog.Height(64) <= 8 {
		t.Errorf("Expected height to grow beyond %d, got %d", 8, dialog.Height(64))
	}

}

func TestQuestionDialogHeightCapsInput(t *testing.T) {

	dialog := NewQuestionDialog(&types.Question{Type: "Ask", Question: "What is your name?"})
	dialog.Input.SetText("a\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl")

	if dialog.Height(64) != 15 {
		t.Errorf("Expected height %d, got %d", 15, dialog.Height(64))
	}

}

func TestQuestionDialogCtrlCClearsInput(t *testing.T) {

	dialog := NewQuestionDialog(&types.Question{Type: "Ask", Question: "What is your name?"})
	dialog.Input.SetText("some answer")

	answer, submit, cancel := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 'c', Ctrl: true})

	if dialog.Input.GetText() != "" {
		t.Errorf("Expected input to be cleared, got %q", dialog.Input.GetText())
	}

	if answer != "" || submit == true || cancel == true {
		t.Errorf("Expected Ctrl+C to only clear the input")
	}

}
