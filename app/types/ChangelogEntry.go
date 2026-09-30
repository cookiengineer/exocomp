package types

import "strings"
import "time"

type ChangelogEntry struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	File        string    `json:"file"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
}

func (entry *ChangelogEntry) PackagePath() string {

	path := strings.TrimPrefix(entry.File, "./")

	if strings.Contains(path, "/") == true {
		return path[:strings.LastIndex(path, "/")]
	}

	return path

}
