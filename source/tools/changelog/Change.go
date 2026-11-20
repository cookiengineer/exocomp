package changelog

func (tool *Changelog) Change(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Change", path, symbol, description)
}
