package tty

type Buffer struct {
	Width  int
	Height int
	cells  []Cell
}

func NewBuffer(width int, height int) *Buffer {

	if width < 1 {
		width = 1
	}

	if height < 1 {
		height = 1
	}

	return &Buffer{
		Width:  width,
		Height: height,
		cells:  make([]Cell, width*height),
	}

}

func (buffer *Buffer) Clear(style Style) {

	for index := range buffer.cells {
		buffer.cells[index] = Cell{Text: " ", Width: 1, Style: style}
	}

}

func (buffer *Buffer) Resize(width int, height int) {

	if width < 1 {
		width = 1
	}

	if height < 1 {
		height = 1
	}

	cells := make([]Cell, width*height)

	for y := 0; y < min(height, buffer.Height); y++ {

		for x := 0; x < min(width, buffer.Width); x++ {
			cells[y*width+x] = buffer.cells[y*buffer.Width+x]
		}

	}

	buffer.Width = width
	buffer.Height = height
	buffer.cells = cells

}

func (buffer *Buffer) Get(x int, y int) Cell {

	if x >= 0 && x < buffer.Width && y >= 0 && y < buffer.Height {
		return buffer.cells[y*buffer.Width+x]
	} else {
		return Cell{}
	}

}

func (buffer *Buffer) Set(x int, y int, text string, style Style) {

	if x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height {
		return
	}

	width := StringWidth(text)

	if width < 1 {
		width = 1
	}

	buffer.cells[y*buffer.Width+x] = Cell{
		Text:  text,
		Width: width,
		Style: style,
	}

	for offset := 1; offset < width && x+offset < buffer.Width; offset++ {
		buffer.cells[y*buffer.Width+x+offset] = Cell{Text: "", Width: 0, Style: style}
	}

}

func (buffer *Buffer) SetRune(x int, y int, value rune, style Style) {

	width := RuneWidth(value)

	if width <= 0 {
		return
	}

	buffer.Set(x, y, string(value), style)

}

func (buffer *Buffer) Fill(rect Rect, value rune, style Style) {

	clipped := rect.Intersect(Rect{X: 0, Y: 0, Width: buffer.Width, Height: buffer.Height})

	for y := clipped.Y; y < clipped.Y+clipped.Height; y++ {

		for x := clipped.X; x < clipped.X+clipped.Width; x++ {
			buffer.SetRune(x, y, value, style)
		}

	}

}

func (buffer *Buffer) WriteString(x int, y int, text string, style Style) int {

	cursor := x
	runes := []rune(text)

	for index := 0; index < len(runes); index++ {

		value := runes[index]
		width := RuneWidth(value)

		if width == 0 {

			// NOTE: Combining marks attach to the previous cell.
			if cursor > x {

				previous := buffer.Get(cursor-1, y)
				previous.Text = previous.Text + string(value)
				buffer.Set(cursor-1, y, previous.Text, previous.Style)

			}

		} else {

			if cursor < buffer.Width {
				buffer.SetRune(cursor, y, value, style)
			}

			cursor += width

		}

	}

	return cursor - x

}

func (buffer *Buffer) Blit(source *Buffer, destination_x int, destination_y int) {

	for y := 0; y < source.Height; y++ {

		for x := 0; x < source.Width; x++ {

			cell := source.Get(x, y)

			if cell.IsEmpty() == false {
				buffer.Set(destination_x+x, destination_y+y, cell.Text, cell.Style)
			}

		}

	}

}

func (buffer *Buffer) Clone() *Buffer {

	return &Buffer{
		Width:  buffer.Width,
		Height: buffer.Height,
		cells:  append([]Cell{}, buffer.cells...),
	}

}
