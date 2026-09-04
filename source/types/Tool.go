package types

import "exocomp/schemas"

type Tool interface {

	Name()                       string
	Call(string, map[string]any) (string, error)

	// XXX: This method can't be put anywhere else, cyclic dependency loop
	GetContentIdentifiers() []string
	GetContent(string)      (any, error)
	HasMethod(string)       bool
	Schemas()               []schemas.Tool

}
