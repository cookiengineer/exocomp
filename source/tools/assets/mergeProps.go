package assets

func mergeProps(left map[string]string, right map[string]string) map[string]string {

	result := make(map[string]string, len(left)+len(right))

	for key, value := range left {
		result[key] = value
	}

	for key, value := range right {

		if value == "" {
			continue
		} else {
			result[key] = value
		}
	}

	return result

}
