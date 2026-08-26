package markdown

func splitTableCells(line string) []string {

	cells := make([]string, 0)
	current := ""
	in_code := false

	for i := 0; i < len(line); i++ {

		ch := line[i]

		if ch == '`' {

			in_code = !in_code
			current += string(ch)

		} else if ch == '\\' && i+1 < len(line) && line[i+1] == '|' && !in_code {

			current += string('|')
			i++

		} else if ch == '|' && !in_code {

			cells = append(cells, current)
			current = ""

		} else {

			current += string(ch)

		}

	}

	cells = append(cells, current)

	return cells

}
