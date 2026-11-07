package components

import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "testing"

func TestFooterEnterNewline(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 'a'})
	footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyEnter})
	footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 'b'})

	if footer.GetText() != "a\nb" {
		t.Errorf("Expected %q, got %q", "a\nb", footer.GetText())
	}

}

func TestFooterCtrlSSend(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 'a'})

	if footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 's', Ctrl: true}) != true {
		t.Errorf("Expected Ctrl+S to submit")
	}

	if footer.GetText() != "a" {
		t.Errorf("Expected %q, got %q", "a", footer.GetText())
	}

}

func TestFooterSlashPalette(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: '/'})

	if footer.Palette.Visible == true {
		if footer.GetText() != "/" {
			t.Errorf("Expected slash command, got %q", footer.GetText())
		}
	} else {
		t.Errorf("Expected palette visible")
	}

}
