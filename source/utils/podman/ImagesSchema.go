package podman

type ImagesSchema struct {
	ID       string   `json:"Id"`
	RepoTags []string `json:"RepoTags"`
}
