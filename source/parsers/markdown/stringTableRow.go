package markdown

func stringTableRow(row []string, lengths []int, alignment []string, whitespace string) string {

	result := "|"

	for c, cell := range row {

		escaped := escapeTableCells(cell)
		padding := lengths[c] - len(cell)

		switch alignment[c] {

		case "left":
			result += " " + escaped + whitespace[:padding+1] + "|"

		case "right":
			result += whitespace[:padding+1] + escaped + " |"

		case "justify", "center":
			left := padding / 2
			right := padding - left
			result += whitespace[:left+1] + escaped + whitespace[:right+1] + "|"

		}

	}

	result += "\n"

	return result

}
