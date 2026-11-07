package components

import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "time"

type QuestionDialog struct {
	Question *types.Question
	Input    *Textarea
	cursorX  int
	cursorY  int
}

func NewQuestionDialog(question *types.Question) *QuestionDialog {

	return &QuestionDialog{
		Question: question,
		Input:    NewTextarea(),
	}

}

func (dialog *QuestionDialog) GetQuestion() string {
	return dialog.Question.Question
}

func (dialog *QuestionDialog) Width(width int) int {
	return min(64, max(20, width-4))
}

func (dialog *QuestionDialog) Height(width int) int {

	question_lines := len(utils_tty.Wrap(dialog.Question.Question, max(8, width-6)))
	lines, _, _ := dialog.Input.Wrap(max(1, width-6))
	input_height := min(10, max(1, len(lines)))

	return max(8, question_lines+input_height+4)

}

func (dialog *QuestionDialog) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

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

	lines, cursor_line, cursor_column := dialog.Input.Wrap(max(1, rect.Width-6))
	textarea_height := min(10, max(1, len(lines)))
	available_height := max(1, rect.Y+rect.Height-row-2)

	if textarea_height > available_height {
		textarea_height = available_height
	}

	if cursor_line < dialog.Input.scroll {
		dialog.Input.scroll = cursor_line
	} else if cursor_line >= dialog.Input.scroll+textarea_height {
		dialog.Input.scroll = cursor_line - textarea_height + 1
	}

	input_rect := utils_tty.Rect{
		X:      rect.X + 2,
		Y:      row,
		Width:  max(1, rect.Width-4),
		Height: textarea_height,
	}

	buffer.Fill(input_rect, ' ', styles.DialogInput)

	for index := 0; index < input_rect.Height; index++ {

		line_index := dialog.Input.scroll + index

		if line_index >= len(lines) {
			break
		}

		buffer.WriteString(input_rect.X+1, input_rect.Y+index, utils_tty.Truncate(lines[line_index].Text, input_rect.Width-2), styles.DialogInput)

	}

	dialog.cursorX = input_rect.X + 1 + cursor_column
	dialog.cursorY = input_rect.Y + (cursor_line - dialog.Input.scroll)

	if dialog.cursorX >= input_rect.X+input_rect.Width {
		dialog.cursorX = input_rect.X + input_rect.Width - 1
	}

	if dialog.cursorY >= input_rect.Y+input_rect.Height {
		dialog.cursorY = input_rect.Y + input_rect.Height - 1
	}

	buffer.WriteString(
		rect.X+2,
		rect.Y+rect.Height-1,
		utils_tty.Truncate("Enter newline · Ctrl+S submit · Esc cancel", rect.Width-4),
		styles.DialogMuted,
	)

}

func (dialog *QuestionDialog) GetCursor() (int, int, bool) {
	return dialog.cursorX, dialog.cursorY, true
}

func (dialog *QuestionDialog) HandleKey(key utils_tty.Key) (string, bool, bool) {

	if key.Type == utils_tty.KeyRune && key.Ctrl == true && key.Rune == 'c' {

		dialog.Input.Clear()

		return "", false, false

	}

	switch key.Type {

	case utils_tty.KeyEnter:
		dialog.Input.InsertNewline()

	case utils_tty.KeyEscape:
		return "", false, true

	case utils_tty.KeyBackspace:
		dialog.Input.Backspace()

	case utils_tty.KeyDelete:
		dialog.Input.Delete()

	case utils_tty.KeyArrowLeft:
		dialog.Input.MoveLeft()

	case utils_tty.KeyArrowRight:
		dialog.Input.MoveRight()

	case utils_tty.KeyArrowUp:
		dialog.Input.MoveUp()

	case utils_tty.KeyArrowDown:
		dialog.Input.MoveDown()

	case utils_tty.KeyHome:
		dialog.Input.Home()

	case utils_tty.KeyEnd:
		dialog.Input.End()

	case utils_tty.KeyPaste:
		dialog.Input.Insert(key.Text)

	case utils_tty.KeyRune:

		if key.Ctrl == true {

			if key.Rune == 's' {
				return dialog.Input.GetText(), true, false
			} else if key.Rune == 'j' {
				dialog.Input.InsertNewline()
			}

		} else {
			dialog.Input.Insert(string(key.Rune))
		}

	}

	return "", false, false

}
