package changelog

func (tool *Changelog) Deprecate(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Deprecate", path, symbol, description)
}
