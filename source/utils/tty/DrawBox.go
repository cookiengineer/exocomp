package tty

func DrawBox(buffer *Buffer, rect Rect, chars BorderChars, style Style) {

	if rect.Width < 2 || rect.Height < 2 {
		return
	}

	for x := rect.X + 1; x < rect.X+rect.Width-1; x++ {

		buffer.Set(x, rect.Y, chars.Horizontal, style)
		buffer.Set(x, rect.Y+rect.Height-1, chars.Horizontal, style)

	}

	for y := rect.Y + 1; y < rect.Y+rect.Height-1; y++ {

		buffer.Set(rect.X, y, chars.Vertical, style)
		buffer.Set(rect.X+rect.Width-1, y, chars.Vertical, style)

	}

	buffer.Set(rect.X, rect.Y, chars.TopLeft, style)
	buffer.Set(rect.X+rect.Width-1, rect.Y, chars.TopRight, style)
	buffer.Set(rect.X, rect.Y+rect.Height-1, chars.BottomLeft, style)
	buffer.Set(rect.X+rect.Width-1, rect.Y+rect.Height-1, chars.BottomRight, style)

}
