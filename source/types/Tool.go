package types

type Tool interface {

	Call(string, map[string]interface{}) (string, error)

	// XXX: This method can't be put anywhere else, cyclic dependency loop
	GetContent(string) (any, error)

}
