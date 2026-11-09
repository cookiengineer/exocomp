package tty

func DrawSplit(buffer *Buffer, rect Rect, style Style) {

	for y := rect.Y; y < rect.Y+rect.Height; y++ {

		buffer.Set(rect.X, y, BorderSplit.Vertical, style)

		if rect.Width > 1 {
			buffer.Set(rect.X+rect.Width-1, y, BorderSplit.Vertical, style)
		}

	}

}
