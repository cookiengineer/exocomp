package markdown

import "strings"

func renderTable(element *Element, document *Document, indent string) string {

	result := ""
	alignment := make([]string, 0)
	thead := make([]string, 0)
	tbody := make([][]string, 0)
	tfoot := make([][]string, 0)

	lines := strings.Split(strings.TrimSpace(element.Text), "\n")

	if len(lines) >= 2 {

		first_line := splitTableCells(lines[0][1 : len(lines[0])-1])

		for f := 0; f < len(first_line); f++ {
			thead = append(thead, strings.TrimSpace(first_line[f]))
		}

		if strings.Contains(lines[1], "---") {

			second_line := splitTableCells(lines[1][1 : len(lines[1])-1])

			for s := 0; s < len(second_line); s++ {

				cell := second_line[s]

				if strings.HasPrefix(cell, ":-") && strings.HasSuffix(cell, "-:") {
					alignment = append(alignment, "justify")
				} else if strings.HasPrefix(cell, ":-") {
					alignment = append(alignment, "left")
				} else if strings.HasSuffix(cell, "-:") {
					alignment = append(alignment, "right")
				} else {
					alignment = append(alignment, "center")
				}

			}

			if len(lines) > 2 {

				fill := "tbody"

				for l := 2; l < len(lines); l++ {

					line := lines[l]

					if strings.Contains(line, "---") {

						fill = "tfoot"

					} else {

						cells := splitTableCells(line[1 : len(line)-1])
						row := make([]string, 0)

						for _, cell := range cells {
							row = append(row, strings.TrimSpace(cell))
						}

						if fill == "tbody" {
							tbody = append(tbody, row)
						} else if fill == "tfoot" {
							tfoot = append(tfoot, row)
						}

					}

				}

			}

		} else {

			for t := 0; t < len(thead); t++ {
				alignment = append(alignment, "center")
			}

			if len(lines) > 1 {

				fill := "tbody"

				for l := 1; l < len(lines); l++ {

					line := lines[l]

					if strings.Contains(line, "---") {

						fill = "tfoot"

					} else {

						cells := splitTableCells(line[1 : len(line)-1])
						row := make([]string, 0)

						for _, cell := range cells {
							row = append(row, strings.TrimSpace(cell))
						}

						if fill == "tbody" {
							tbody = append(tbody, row)
						} else if fill == "tfoot" {
							tfoot = append(tbody, row)
						}

					}

				}

			}

		}

	}

	if len(thead) > 0 || len(tbody) > 0 || len(tfoot) > 0 {

		result += indent + "<" + element.Type + ">\n"

		if len(thead) > 0 {

			result += indent + "\t<thead>\n"
			result += indent + "\t\t<tr>\n"

			for t := 0; t < len(thead); t++ {

				elements, err := parseInlineElements(strings.TrimSpace(thead[t]), element.Line, false)

				if err == nil {

					align := alignment[t]

					if align == "center" {
						result += indent + "\t\t\t<th align=\"center\">"
					} else if align == "left" {
						result += indent + "\t\t\t<th align=\"left\">"
					} else if align == "right" {
						result += indent + "\t\t\t<th align=\"right\">"
					} else if align == "justify" {
						result += indent + "\t\t\t<th align=\"justify\">"
					}

					for e := 0; e < len(elements); e++ {
						result += elements[e].RenderInto(document, "")
					}

					result += "</th>\n"

				}

			}

			result += indent + "\t\t</tr>\n"
			result += indent + "\t</thead>\n"

		}

		if len(tbody) > 0 {

			result += indent + "\t<tbody>\n"

			for t := 0; t < len(tbody); t++ {

				result += indent + "\t\t<tr>\n"

				cells := tbody[t]

				for c := 0; c < len(cells); c++ {

					elements, err := parseInlineElements(strings.TrimSpace(cells[c]), element.Line+1+t+1, false)

					if err == nil {

						align := alignment[c]

						if align == "center" {
							result += indent + "\t\t\t<td align=\"center\">"
						} else if align == "left" {
							result += indent + "\t\t\t<td align=\"left\">"
						} else if align == "right" {
							result += indent + "\t\t\t<td align=\"right\">"
						} else if align == "justify" {
							result += indent + "\t\t\t<td align=\"justify\">"
						}

						for e := 0; e < len(elements); e++ {
							result += elements[e].RenderInto(document, "")
						}

						result += "</td>\n"

					}

				}

				result += indent + "\t\t</tr>\n"

			}

			result += indent + "\t</tbody>\n"

		}

		if len(tfoot) > 0 {

			result += indent + "\t<tfoot>\n"

			for t := 0; t < len(tfoot); t++ {

				cells := tfoot[t]

				result += indent + "\t\t<tr>\n"

				for c := 0; c < len(cells); c++ {

					elements, err := parseInlineElements(strings.TrimSpace(cells[c]), element.Line+1+len(tbody)+t+1, false)

					if err == nil {

						align := alignment[c]

						if align == "center" {
							result += indent + "\t\t\t<td align=\"center\">"
						} else if align == "left" {
							result += indent + "\t\t\t<td align=\"left\">"
						} else if align == "right" {
							result += indent + "\t\t\t<td align=\"right\">"
						} else if align == "justify" {
							result += indent + "\t\t\t<td align=\"justify\">"
						}

						for e := 0; e < len(elements); e++ {
							result += elements[e].RenderInto(document, "")
						}

						result += "</td>\n"

					}

				}

				result += indent + "\t\t</tr>\n"

			}

			result += indent + "\t</tfoot>\n"

		}

		result += indent + "</" + element.Type + ">"

	}

	return result

}
