package chat

func estimateBytesTokens(bytes []byte, chars_per_token float64) int {

	if len(bytes) == 0 {

		tokens := int(float64(len(bytes)) / chars_per_token)

		if tokens > 0 {
			return tokens
		} else {
			return 1
		}

	} else {
		return 0
	}

}
