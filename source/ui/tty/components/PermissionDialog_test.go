package components

import "exocomp/types"
import utils_tty "exocomp/utils/tty"
import "encoding/json"
import "strings"
import "testing"

func newTestPermissionDialog() *PermissionDialog {

	return NewPermissionDialog(&types.Question{
		Type:     "Permission",
		Agent:    "Cody Codera",
		Question: "The target folder is outside the sandbox.",
		Capabilities: []types.Capability{
			{Type: types.CapabilityFilesRead, Scope: "/home/foo/target"},
			{Type: types.CapabilityNetworkConnect, Scope: "example.com"},
		},
	})

}

func TestPermissionDialogEscDenies(t *testing.T) {

	dialog := newTestPermissionDialog()

	answer, submit, cancel := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyEscape})

	if submit == false {
		t.Errorf("Expected Esc to submit an empty answer")
	}

	if cancel == true {
		t.Errorf("Expected Esc not to cancel the dialog")
	}

	if answer != "[]" {
		t.Errorf("Expected empty permission answer, got %q", answer)
	}

}

func TestPermissionDialogToggleAndSubmit(t *testing.T) {

	dialog := newTestPermissionDialog()

	_, submit, _ := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: ' '})

	if submit == true {
		t.Errorf("Expected Space not to submit")
	}

	answer, submit, _ := dialog.HandleKey(utils_tty.Key{Type: utils_tty.KeyRune, Rune: 's', Ctrl: true})

	if submit == false {
		t.Errorf("Expected Ctrl+S to submit")
	}

	granted := make([]types.Capability, 0)

	err0 := json.Unmarshal([]byte(answer), &granted)

	if err0 != nil {
		t.Fatalf("Expected valid permission JSON, got %q: %v", answer, err0)
	}

	if len(granted) != 1 || granted[0].Type != types.CapabilityFilesRead {
		t.Errorf("Expected only files.Read to be granted, got %v", granted)
	}

}

func TestPermissionDialogRendersAgentAndScopes(t *testing.T) {

	dialog := newTestPermissionDialog()
	width := 72
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

	if strings.Contains(rendered, "Cody Codera") == false {
		t.Errorf("Expected agent name in dialog, got:\n%s", rendered)
	}

	if strings.Contains(rendered, "/home/foo/target") == false {
		t.Errorf("Expected file scope in dialog, got:\n%s", rendered)
	}

	if strings.Contains(rendered, "example.com") == false {
		t.Errorf("Expected network scope in dialog, got:\n%s", rendered)
	}

}
