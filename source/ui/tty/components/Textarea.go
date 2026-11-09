package components

import "strings"

type Textarea struct {
	runes  []rune
	cursor int
	scroll int
}

func NewTextarea() *Textarea {

	return &Textarea{
		runes:  make([]rune, 0),
		cursor: 0,
		scroll: 0,
	}

}

func (textarea *Textarea) Clear() {

	textarea.runes = make([]rune, 0)
	textarea.cursor = 0
	textarea.scroll = 0

}

func (textarea *Textarea) GetCursor() int {
	return textarea.cursor
}

func (textarea *Textarea) GetLineStart() (int, int) {

	start := textarea.cursor

	for start > 0 && textarea.runes[start-1] != '\n' {
		start--
	}

	return start, textarea.cursor - start

}

func (textarea *Textarea) GetText() string {
	return string(textarea.runes)
}

func (textarea *Textarea) Insert(value string) {

	runes := []rune(value)

	before := append([]rune{}, textarea.runes[0:textarea.cursor]...)
	after := append([]rune{}, textarea.runes[textarea.cursor:]...)

	textarea.runes = append(before, runes...)
	textarea.runes = append(textarea.runes, after...)
	textarea.cursor += len(runes)

}

func (textarea *Textarea) InsertNewline() {
	textarea.Insert("\n")
}
func (textarea *Textarea) SetText(value string) {

	textarea.runes = []rune(value)
	textarea.cursor = len(textarea.runes)
	textarea.scroll = 0

}

func (textarea *Textarea) Delete() bool {

	if textarea.cursor < len(textarea.runes) {

		textarea.runes = append(textarea.runes[0:textarea.cursor], textarea.runes[textarea.cursor+1:]...)

		return true

	} else {
		return false
	}

}

func (textarea *Textarea) MoveLeft() {

	if textarea.cursor > 0 {
		textarea.cursor--
	}

}

func (textarea *Textarea) MoveRight() {

	if textarea.cursor < len(textarea.runes) {
		textarea.cursor++
	}

}

func (textarea *Textarea) MoveUp() {

	start, column := textarea.GetLineStart()
	_ = column

	if start > 0 {

		previous_end := start - 1
		previous_start := previous_end

		for previous_start > 0 && textarea.runes[previous_start-1] != '\n' {
			previous_start--
		}

		target := previous_start + min(textarea.cursor-start, previous_end-previous_start)
		textarea.cursor = target

	}

}

func (textarea *Textarea) MoveDown() {

	start, column := textarea.GetLineStart()
	_ = column

	end := start

	for end < len(textarea.runes) && textarea.runes[end] != '\n' {
		end++
	}

	if end < len(textarea.runes) {

		next_start := end + 1
		next_end := next_start

		for next_end < len(textarea.runes) && textarea.runes[next_end] != '\n' {
			next_end++
		}

		textarea.cursor = next_start + min(textarea.cursor-start, next_end-next_start)

	}

}

func (textarea *Textarea) Backspace() bool {

	if textarea.cursor > 0 {

		textarea.runes = append(textarea.runes[0:textarea.cursor-1], textarea.runes[textarea.cursor:]...)
		textarea.cursor--

		return true

	} else {
		return false
	}

}

func (textarea *Textarea) Home() {

	start, _ := textarea.GetLineStart()
	textarea.cursor = start

}

func (textarea *Textarea) End() {

	_, _ = textarea.GetLineStart()
	end := textarea.cursor

	for end < len(textarea.runes) && textarea.runes[end] != '\n' {
		end++
	}

	textarea.cursor = end

}

func (textarea *Textarea) Wrap(width int) ([]TextareaLine, int, int) {

	if width <= 0 {
		width = 1
	}

	lines := make([]TextareaLine, 0)
	builder := strings.Builder{}
	line_width := 0
	cursor_line := 0
	cursor_column := 0
	found_cursor := false

	for index := 0; index <= len(textarea.runes); index++ {

		if index == textarea.cursor {

			cursor_line = len(lines)
			cursor_column = line_width
			found_cursor = true

		}

		if index == len(textarea.runes) {
			break
		}

		value := textarea.runes[index]

		if value == '\n' {

			lines = append(lines, TextareaLine{
				Text:         builder.String(),
				CursorColumn: -1,
				CursorHere:   false,
			})
			builder.Reset()
			line_width = 0
			continue

		}

		value_width := runeWidth(value)

		if line_width+value_width > width && line_width > 0 {

			lines = append(lines, TextareaLine{
				Text:         builder.String(),
				CursorColumn: -1,
				CursorHere:   false,
			})
			builder.Reset()
			line_width = 0

		}

		builder.WriteRune(value)
		line_width += value_width

	}

	lines = append(lines, TextareaLine{
		Text:         builder.String(),
		CursorColumn: -1,
		CursorHere:   false,
	})

	if found_cursor == false {

		cursor_line = len(lines) - 1
		cursor_column = line_width

	}

	if len(lines) == 0 {
		lines = append(lines, TextareaLine{Text: ""})
	}

	if cursor_line >= len(lines) {
		cursor_line = len(lines) - 1
	}

	if cursor_line >= 0 && cursor_line < len(lines) {

		lines[cursor_line].CursorColumn = cursor_column
		lines[cursor_line].CursorHere = true

	}

	return lines, cursor_line, cursor_column

}
