package components

import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "strings"
import "time"

type ChoiceDialog struct {
	Question *types.Question
	Selected int
	chosen   map[int]bool
}

func NewChoiceDialog(question *types.Question) *ChoiceDialog {

	return &ChoiceDialog{
		Question: question,
		Selected: 0,
		chosen:   make(map[int]bool),
	}

}

func (dialog *ChoiceDialog) GetQuestion() string {
	return dialog.Question.Question
}

func (dialog *ChoiceDialog) Width(width int) int {
	return min(64, max(20, width-4))
}

func (dialog *ChoiceDialog) Height(width int) int {

	question_lines := len(utils_tty.Wrap(dialog.Question.Question, max(8, width-6)))

	return question_lines + len(dialog.Question.Options) + 5

}

func (dialog *ChoiceDialog) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 2 || rect.Height <= 2 {
		return
	}

	buffer.Fill(rect, ' ', styles.Dialog)
	utils_tty.DrawBox(buffer, rect, utils_tty.BorderRounded, styles.DialogBorder)
	buffer.WriteString(rect.X+2, rect.Y, " Human Question ", styles.DialogTitle)

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

	for index, option := range dialog.Question.Options {

		if row >= rect.Y+rect.Height-1 {
			break
		}

		selected := index == dialog.Selected
		marker := "[ ]"

		if dialog.Question.Multiple == true {

			if dialog.chosen[index] == true {
				marker = "[x]"
			}

		} else if selected == true {
			marker = "[x]"
		}

		style := styles.Dialog

		if selected == true {

			style = styles.Selection
			buffer.Fill(utils_tty.Rect{X: rect.X + 2, Y: row, Width: rect.Width - 4, Height: 1}, ' ', style)

		}

		text := marker + " " + option
		buffer.WriteString(rect.X+2, row, utils_tty.Truncate(text, rect.Width-4), style)
		row++

	}

	hint := "↑/↓ select · Ctrl+S submit · Esc cancel"

	if dialog.Question.Multiple == true {
		hint = "↑/↓ select · Space toggle · Ctrl+S submit · Esc cancel"
	}

	buffer.WriteString(
		rect.X+2,
		rect.Y+rect.Height-1,
		utils_tty.Truncate(hint, rect.Width-4),
		styles.DialogMuted,
	)

}

func (dialog *ChoiceDialog) GetCursor() (int, int, bool) {
	return 0, 0, false
}

func (dialog *ChoiceDialog) HandleKey(key utils_tty.Key) (string, bool, bool) {

	if key.Type == utils_tty.KeyRune && key.Ctrl == true && key.Rune == 'c' {

		dialog.Selected = 0
		dialog.chosen = make(map[int]bool)

		return "", false, false

	}

	switch key.Type {

	case utils_tty.KeyArrowUp:

		dialog.Selected--

		if dialog.Selected < 0 {
			dialog.Selected = len(dialog.Question.Options) - 1
		}

	case utils_tty.KeyArrowDown:

		dialog.Selected++

		if dialog.Selected >= len(dialog.Question.Options) {
			dialog.Selected = 0
		}

	case utils_tty.KeyRune:

		if key.Ctrl == true && key.Rune == 's' {

			answers := make([]string, 0)

			if dialog.Question.Multiple == true {

				for index, option := range dialog.Question.Options {

					if dialog.chosen[index] == true {
						answers = append(answers, option)
					}

				}

			} else if dialog.Selected >= 0 && dialog.Selected < len(dialog.Question.Options) {
				answers = append(answers, dialog.Question.Options[dialog.Selected])
			}

			return strings.Join(answers, ", "), true, false

		}

		if key.Rune == ' ' {

			if dialog.Question.Multiple == true {
				dialog.chosen[dialog.Selected] = !dialog.chosen[dialog.Selected]
			}

		}

	case utils_tty.KeyEscape:
		return "", false, true

	}

	return "", false, false

}
