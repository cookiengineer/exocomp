package findings

// NewAssetID returns a new random asset identifier.
func NewAssetID() string {
	return newPrefixedID("ast_")
}
