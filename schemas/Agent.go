package schemas

type Agent struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Model       string     `json:"model"`
	Temperature float64    `json:"temperature"`
	Programs    []string   `json:"programs"`
	Tools       []string   `json:"tools"`
	Messages    []*Message `json:"messages"`
	Sandbox     string     `json:"sandbox"`
}
