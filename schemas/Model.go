package schemas

import utils_time "exocomp/utils/time"

type Model struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created utils_time.UnixTime `json:"created"`
	OwnedBy string              `json:"owned_by"`
}
