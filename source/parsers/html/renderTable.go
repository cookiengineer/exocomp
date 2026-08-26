package html

import "strings"

func renderTable(table *Element, document *Document) string {

	rows := make([][]string, 0)

	for _, child := range table.Children {

		if child.Type == "tr" {

			cells := make([]string, 0)

			for _, cell := range child.Children {

				if cell.Type == "th" || cell.Type == "td" {
					cells = append(cells, escapeTableCell(cell.textContent()))
				}

			}

			rows = append(rows, cells)

		}

	}

	if len(rows) == 0 {
		return ""
	}

	max_cols := 0

	for _, row := range rows {

		if len(row) > max_cols {
			max_cols = len(row)
		}

	}

	for r := 0; r < len(rows); r++ {

		for len(rows[r]) < max_cols {
			rows[r] = append(rows[r], "")
		}

	}

	lines := make([]string, 0)
	lines = append(lines, "| "+strings.Join(rows[0], " | ")+" |")

	separator := make([]string, 0)

	for c := 0; c < max_cols; c++ {
		separator = append(separator, "---")
	}

	lines = append(lines, "| "+strings.Join(separator, " | ")+" |")

	for r := 1; r < len(rows); r++ {
		lines = append(lines, "| "+strings.Join(rows[r], " | ")+" |")
	}

	return strings.Join(lines, "\n")

}

func escapeTableCell(value string) string {

	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")

	return strings.TrimSpace(value)

}
