package build

type Program struct {
	Name       string `json:"name"`
	Folder     string `json:"folder"`
	Repository string `json:"repository"`
	Version    string `json:"version"`
	Source     string `json:"source"`
}
