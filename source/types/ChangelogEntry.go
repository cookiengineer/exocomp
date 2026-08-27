package types

import "time"

type ChangelogEntry struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	File        string    `json:"file"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
}
