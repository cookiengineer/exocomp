package types

import "exocomp/schemas"

type Tool interface {

	Name()                               string
	Call(string, map[string]interface{}) (string, error)

	// XXX: This method can't be put anywhere else, cyclic dependency loop
	GetContent(string) (any, error)
	Schemas()          []schemas.Tool

}
