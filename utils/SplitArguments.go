package utils

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isPrintableASCII(b byte) bool {
	return b >= 0x20 && b <= 0x7E
}

func SplitArguments(input string) []string {

	args    := make([]string, 0)
	current := make([]byte, 0)

	in_double_quotes := false
	in_single_quotes := false
	escaped          := false

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if isPrintableASCII(chr) == true {

			switch {
			case escaped:

				if in_single_quotes == true {

					escaped = false

				} else {

					if chr == '"' || chr == '\\' || isSpace(chr) {
						current = append(current, chr)
						escaped = false
					}

				}

			case chr == '\\':

				if in_single_quotes {
					current = append(current, chr)
				} else {
					escaped = true
				}

			case chr == '"' && in_single_quotes == false:
				in_double_quotes = !in_double_quotes

			case chr == '\'' && in_double_quotes == false:
				in_single_quotes = !in_single_quotes

			case isSpace(chr) && in_single_quotes == false && in_double_quotes == false:

				if len(current) > 0 {
					args    = append(args, string(current))
					current = current[0:0]
				}

			default:
				current = append(current, chr)

			}

		}

	}

	if len(current) > 0 {
		args = append(args, string(current))
	}

	return args
}

