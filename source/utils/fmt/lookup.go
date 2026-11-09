package fmt

func lookup(character uint16) string {

	var found string

	for str, codes := range utf8Table {

		for c := 0; c < len(codes); c++ {

			if codes[c] == character {

				found = str
				break

			}

		}

	}

	return found

}
