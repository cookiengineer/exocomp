package strings

func IsNumber(value string) bool {

	var result bool = true

	for v := 0; v < len(value); v++ {

		character := string(value[v])

		if character >= "0" && character <= "9" {
			continue
		} else if character == "+" || character == "-" || character == "." {
			continue
		} else if character == "E" || character == "e" {
			continue
		} else {
			result = false
			break
		}

	}

	return result

}
