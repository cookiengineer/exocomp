package types

import "strings"

type Bug struct {
	IsFixed     bool   `json:"is_fixed"`
	File        string `json:"file"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

func (bug *Bug) PackagePath() string {

	path := strings.TrimPrefix(bug.File, "./")

	if strings.Contains(path, "/") == true {
		return path[:strings.LastIndex(path, "/")]
	}

	return path

}
