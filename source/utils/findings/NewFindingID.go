package findings

func NewFindingID() string {
	return newPrefixedID("fnd_")
}
