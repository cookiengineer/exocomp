package tools

// shellArgumentBytes returns the combined byte length of all arguments.
func shellArgumentBytes(arguments []string) int {

	size := 0

	for _, argument := range arguments {
		size += len(argument)
	}

	return size

}
