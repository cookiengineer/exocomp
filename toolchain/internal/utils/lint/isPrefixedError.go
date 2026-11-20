package lint

func isPrefixedError(message string) bool {

	if prefixedErrorPattern.MatchString(message) == true {
		return true
	} else {
		return parameterPrefixPattern.MatchString(message)
	}

}
