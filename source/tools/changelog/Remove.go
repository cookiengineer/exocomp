package changelog

func (tool *Changelog) Remove(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Remove", path, symbol, description)
}
