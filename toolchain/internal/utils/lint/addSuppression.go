package lint

func addSuppression(target map[int]map[string]bool, line int, rule string) {

	if target[line] == nil {
		target[line] = make(map[string]bool)
	}

	target[line][rule] = true

}
