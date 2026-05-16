package llamacpp

import utils_time "exocomp/utils/time"

type Model struct {

	ID      string              `json:"id"`
	Aliases []string            `json:"aliases"`
	Object  string              `json:"object"`
	Created utils_time.UnixTime `json:"created"`
	OwnedBy string              `json:"owned_by"`
	Meta    ModelMeta           `json:"meta"`

}
