package tty

func DrawRule(buffer *Buffer, x int, y int, width int, style Style) {

	if width > 0 {
		buffer.WriteString(x, y, repeatText(BorderLight.Horizontal, width), style)
	}

}
