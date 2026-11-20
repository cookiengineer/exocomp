package lint

func camelSegments(name string) []string {

	runes := []rune(name)
	segments := make([]string, 0)
	current := make([]rune, 0)

	for index, character := range runes {

		if isUpperRune(character) == true && len(current) > 0 {

			previous_upper := isUpperRune(runes[index-1])
			next_lower := index+1 < len(runes) && isLowerRune(runes[index+1])

			if previous_upper == false || next_lower == true {

				segments = append(segments, string(current))
				current = make([]rune, 0)

			}

		}

		current = append(current, character)

	}

	if len(current) > 0 {
		segments = append(segments, string(current))
	}

	return segments

}
