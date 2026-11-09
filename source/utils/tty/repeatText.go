package tty

func repeatText(text string, count int) string {

	result := make([]byte, 0)

	for count > 0 {

		result = append(result, text...)
		count--

	}

	return string(result)

}
