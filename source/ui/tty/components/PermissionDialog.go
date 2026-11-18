package components

import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "encoding/json"
import "strings"
import "time"

type PermissionDialog struct {
	Question *types.Question
	Selected int
	chosen   map[int]bool
}

func NewPermissionDialog(question *types.Question) *PermissionDialog {

	return &PermissionDialog{
		Question: question,
		Selected: 0,
		chosen:   make(map[int]bool),
	}

}

func (dialog *PermissionDialog) GetQuestion() string {
	return dialog.Question.Question
}

func (dialog *PermissionDialog) Width(width int) int {
	return min(72, max(24, width-4))
}

func (dialog *PermissionDialog) Height(width int) int {

	question_lines := len(utils_tty.Wrap(dialog.Question.Question, max(8, width-6)))

	return question_lines + len(dialog.Question.Capabilities) + 6

}

func (dialog *PermissionDialog) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 2 || rect.Height <= 2 {
		return
	}

	buffer.Fill(rect, ' ', styles.Dialog)
	utils_tty.DrawBox(buffer, rect, utils_tty.BorderRounded, styles.DialogBorder)

	title := " Agent requests permissions "

	if strings.TrimSpace(dialog.Question.Agent) != "" {
		title = " Agent " + strings.TrimSpace(dialog.Question.Agent) + " requests permissions "
	}

	buffer.WriteString(rect.X+2, rect.Y, utils_tty.Truncate(title, rect.Width-4), styles.DialogTitle)

	question_lines := utils_tty.Wrap(dialog.Question.Question, max(8, rect.Width-6))
	row := rect.Y + 2

	for _, line := range question_lines {

		if row >= rect.Y+rect.Height-1 {
			break
		}

		buffer.WriteString(rect.X+2, row, utils_tty.Truncate(line, rect.Width-4), styles.Dialog)
		row++

	}

	row++

	for index, capability := range dialog.Question.Capabilities {

		if row >= rect.Y+rect.Height-1 {
			break
		}

		selected := index == dialog.Selected
		marker := "[ ]"

		if dialog.chosen[index] == true {
			marker = "[x]"
		}

		style := styles.Dialog

		if selected == true {

			style = styles.Selection
			buffer.Fill(utils_tty.Rect{X: rect.X + 2, Y: row, Width: rect.Width - 4, Height: 1}, ' ', style)

		}

		text := marker + " " + capability.Type + "  " + capability.Scope
		buffer.WriteString(rect.X+2, row, utils_tty.Truncate(text, rect.Width-4), style)
		row++

	}

	buffer.WriteString(
		rect.X+2,
		rect.Y+rect.Height-1,
		utils_tty.Truncate("↑/↓ select · Space toggle · Ctrl+S allow · Esc deny", rect.Width-4),
		styles.DialogMuted,
	)

	if containsNetworkCapability(dialog.Question.Capabilities) == true && rect.Height >= 8 {
		buffer.WriteString(
			rect.X+2,
			rect.Y+rect.Height-2,
			utils_tty.Truncate("Network grants apply to Exocomp network tools only.", rect.Width-4),
			styles.DialogMuted,
		)
	}

}

func (dialog *PermissionDialog) GetCursor() (int, int, bool) {
	return 0, 0, false
}

func (dialog *PermissionDialog) HandleKey(key utils_tty.Key) (string, bool, bool) {

	switch key.Type {

	case utils_tty.KeyArrowUp:

		dialog.Selected--

		if dialog.Selected < 0 {
			dialog.Selected = len(dialog.Question.Capabilities) - 1
		}

	case utils_tty.KeyArrowDown:

		dialog.Selected++

		if dialog.Selected >= len(dialog.Question.Capabilities) {
			dialog.Selected = 0
		}

	case utils_tty.KeyRune:

		if key.Ctrl == true && key.Rune == 's' {
			return dialog.answer(), true, false
		}

		if key.Rune == ' ' && len(dialog.Question.Capabilities) > 0 {
			dialog.chosen[dialog.Selected] = dialog.chosen[dialog.Selected] == false
		}

	case utils_tty.KeyEscape:
		return "[]", true, false

	}

	return "", false, false

}

func (dialog *PermissionDialog) answer() string {

	granted := make([]types.Capability, 0)

	for index, capability := range dialog.Question.Capabilities {

		if dialog.chosen[index] == true {
			granted = append(granted, capability)
		}

	}

	result, err0 := json.Marshal(granted)

	if err0 == nil {
		return string(result)
	} else {
		return "[]"
	}

}
