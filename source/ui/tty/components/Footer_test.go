package components

import "exocomp/ui/tty/styles"
import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "strings"
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

func TestFooterRenderPromptPanel(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	buffer := utils_tty.NewBuffer(80, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 80, Height: 4}, 0)

	if buffer.Get(0, 0).Text != "┃" {
		t.Errorf("Expected the prompt left rail, got %q", buffer.Get(0, 0).Text)
	}

	if buffer.Get(0, 2).Text != "╹" {
		t.Errorf("Expected the prompt rail taper, got %q", buffer.Get(0, 2).Text)
	}

	if buffer.Get(1, 2).Text != "▀" {
		t.Errorf("Expected the prompt bottom cap, got %q", buffer.Get(1, 2).Text)
	}

	if buffer.Get(1, 0).Style != styles.PromptPanel {
		t.Errorf("Expected the composer background")
	}

	if buffer.Get(3, 1).Text != "A" {
		t.Errorf("Expected the placeholder to start after the rail, got %q", buffer.Get(3, 1).Text)
	}

	if buffer.Get(0, 3).Style != styles.Background {
		t.Errorf("Expected the footer row outside the composer")
	}

	for y := 0; y < 4; y++ {

		if strings.Contains(renderTestLine(buffer, y, 80), "─") == true {
			t.Errorf("Expected no composer rule at row %d", y)
		}

	}

}

func TestFooterRenderStatusAndHints(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.Status = "Waiting for the LLM response..."
	buffer := utils_tty.NewBuffer(80, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 80, Height: 4}, 0)

	status_line := renderTestLine(buffer, 3, 80)

	if strings.Contains(status_line, "Agent is working") == false {
		t.Errorf("Expected the working status, got %q", status_line)
	}

	if strings.Contains(status_line, "Ctrl+↑↓ agents") == false {
		t.Errorf("Expected the agent binding, got %q", status_line)
	}

	if strings.Contains(status_line, "Alt+↑↓ scroll") == false {
		t.Errorf("Expected the scroll binding, got %q", status_line)
	}

	if strings.Index(status_line, "Alt+↑↓ scroll") > strings.Index(status_line, "Ctrl+Q quit") {
		t.Errorf("Expected scroll before quit, got %q", status_line)
	}

}

func TestFooterRenderIdleInstruction(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	buffer := utils_tty.NewBuffer(120, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 120, Height: 4}, 0)

	status_line := renderTestLine(buffer, 3, 120)

	if strings.Contains(status_line, "Waiting for") == false {
		t.Errorf("Expected the idle instruction, got %q", status_line)
	}

	if strings.Contains(status_line, "Ctrl+↑↓ agents") == false {
		t.Errorf("Expected the agent binding, got %q", status_line)
	}

	if strings.Contains(status_line, "Alt+↑↓ scroll") == false {
		t.Errorf("Expected the scroll binding, got %q", status_line)
	}

}

func TestFooterRenderBindingKeyStyles(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	buffer := utils_tty.NewBuffer(120, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 120, Height: 4}, 0)

	bindings := footerBindings("")
	x := 118 - footerBindingsWidth(bindings)

	for index, binding := range bindings {

		if buffer.Get(x, 3).Style != styles.FooterKey {
			t.Errorf("Expected the %s key to use the key style", binding.Key)
		}

		x += utils_tty.StringWidth(binding.Key) + 1

		if buffer.Get(x, 3).Style != styles.FooterHint {
			t.Errorf("Expected the %s label to use the hint style", binding.Key)
		}

		x += utils_tty.StringWidth(binding.Label)

		if index < len(bindings)-1 {

			if buffer.Get(x, 3).Style != styles.FooterHint {
				t.Errorf("Expected the separator after %s to use the hint style", binding.Key)
			}

			x += utils_tty.StringWidth(" · ")

		}

	}

}

func TestFooterRenderBusyBindingStyles(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.Status = "Waiting for the LLM response..."
	buffer := utils_tty.NewBuffer(120, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 120, Height: 4}, 0)

	bindings := footerBindings("busy")
	x := 118 - footerBindingsWidth(bindings)

	for index, binding := range bindings {

		if buffer.Get(x, 3).Style != styles.FooterKey {
			t.Errorf("Expected the %s key to use the key style", binding.Key)
		}

		x += utils_tty.StringWidth(binding.Key) + 1

		if buffer.Get(x, 3).Style != styles.FooterHint {
			t.Errorf("Expected the %s label to use the hint style", binding.Key)
		}

		x += utils_tty.StringWidth(binding.Label)

		if index < len(bindings)-1 {
			x += utils_tty.StringWidth(" · ")
		}

	}

}

func TestFooterRenderMultilineErrorStatus(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.Status = "Session.InferChatCompletions: failed\nconnection refused"
	buffer := utils_tty.NewBuffer(120, 4)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 120, Height: 4}, 0)

	status_line := renderTestLine(buffer, 3, 120)

	if strings.Contains(status_line, "connection refused") == false {
		t.Errorf("Expected the error text, got %q", status_line)
	}

	if strings.Contains(status_line, "\n") == true {
		t.Errorf("Expected a single line status, got %q", status_line)
	}

	if strings.Index(status_line, "Alt+↑↓ scroll") > strings.Index(status_line, "Ctrl+Q quit") {
		t.Errorf("Expected scroll before quit, got %q", status_line)
	}

}

func TestFooterRenderMultilineFitsHeight(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.Textarea.SetText("first\nsecond\nthird")
	height := footer.Height(80)
	buffer := utils_tty.NewBuffer(80, height)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 80, Height: height}, 0)

	if buffer.Get(0, height-2).Text != "╹" {
		t.Errorf("Expected the cap on the second last row, got %q", buffer.Get(0, height-2).Text)
	}

	if buffer.Get(0, height-1).Style != styles.Background {
		t.Errorf("Expected the footer on the last row")
	}

}

func TestFooterPaletteStaysInsidePanel(t *testing.T) {

	footer := NewFooter(commands.NewCatalog(nil))
	footer.Palette.Set(newPaletteTestItems(20))

	for index := 0; index < 12; index++ {
		footer.Palette.Move(1)
	}

	buffer := utils_tty.NewBuffer(80, 8)
	footer.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 80, Height: 8}, 0)

	if footer.Palette.scroll == 0 {
		t.Errorf("Expected the palette to scroll")
	}

	if footer.Palette.Selected < footer.Palette.scroll || footer.Palette.Selected >= footer.Palette.scroll+4 {
		t.Errorf("Expected the selected item inside the visible palette window")
	}

	if buffer.Get(4, 5).Style != styles.Selection {
		t.Errorf("Expected the selected item to be visible inside the palette")
	}

	if buffer.Get(1, 6).Text != "▀" {
		t.Errorf("Expected the cap row to stay intact, got %q", buffer.Get(1, 6).Text)
	}

	if buffer.Get(0, 7).Style != styles.Background {
		t.Errorf("Expected the status row to stay outside the palette")
	}

}
