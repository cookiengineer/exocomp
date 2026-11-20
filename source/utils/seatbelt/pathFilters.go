package seatbelt

func pathFilters(paths []resolvedPath) string {

	result := ""

	for index, resolved_path := range paths {

		if index > 0 {
			result = result + " "
		}

		result = result + pathFilter(resolved_path)

	}

	return result

}
