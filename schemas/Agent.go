package schemas

type Agent struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Model       string     `json:"model"`
	Programs    []string   `json:"programs"`
	Temperature float64    `json:"temperature"`
	Tools       []string   `json:"tools"`
	Messages    []*Message `json:"messages"`
}
