package lint

func isExported(name string) bool {

	if name == "" {
		return false
	} else {
		return name[0] >= 'A' && name[0] <= 'Z'
	}

}
