package lint

func groupName(group int) string {

	switch group {

	case groupProject:
		return "project"

	case groupStdlib:
		return "standard library"

	default:
		return "third-party"

	}

}
