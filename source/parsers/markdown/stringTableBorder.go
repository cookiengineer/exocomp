package markdown

func stringTableBorder(lengths []int, alignment []string, border string) string {

	result := "|"

	for c, length := range lengths {

		switch alignment[c] {

		case "left":
			result += ":" + border[1:length+2] + "|"

		case "right":
			result += border[:length+1] + ":|"

		case "justify":
			result += ":" + border[1:length+1] + ":|"

		default:
			result += border[:length+2] + "|"

		}

	}

	result += "\n"

	return result

}
