package changelog

func (tool *Changelog) Fix(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Fix", path, symbol, description)
}
