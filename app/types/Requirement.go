package types

import "strings"

type Requirement struct {
	Type          string `json:"type"`
	File          string `json:"file"`
	Symbol        string `json:"symbol"`
	Declaration   string `json:"declaration"`
	Behavior      string `json:"behavior"`
	IsImplemented bool   `json:"is_implemented"`
}

func (requirement *Requirement) PackagePath() string {

	path := strings.TrimPrefix(requirement.File, "./")

	if strings.Contains(path, "/") == true {
		return path[:strings.LastIndex(path, "/")]
	}

	return path

}
