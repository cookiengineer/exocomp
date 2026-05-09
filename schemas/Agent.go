package schemas

type Agent struct {
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Model           string     `json:"model"`
	Temperature     float64    `json:"temperature"`
	Messages        []*Message `json:"messages"`
	AllowedPrograms []string   `json:"allowed-programs"`
	AllowedTools    []string   `json:"allowed-tools"`
	Sandbox         string     `json:"sandbox"`
}
