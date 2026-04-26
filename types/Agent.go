package types

import "exocomp/schemas"

type Agent struct {
	Name        string
	Type        string
	Model       string
	Prompt      string
	Temperature float64
	Messages    []*schemas.Message
	Programs    []string
	Tools       []string
	Sandbox     string
}

