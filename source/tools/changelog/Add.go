package changelog

func (tool *Changelog) Add(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Add", path, symbol, description)
}
