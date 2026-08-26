package markdown

func escapeTableCells(cell string) string {

	result := ""
	in_code := false

	for i := 0; i < len(cell); i++ {

		ch := cell[i]

		if ch == '`' {
			in_code = !in_code
			result += string(ch)
		} else if ch == '|' && !in_code {
			result += "\\|"
		} else {
			result += string(ch)
		}

	}

	return result

}
