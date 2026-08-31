package schemas

type Error struct {
	Wrapped struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   any     `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

func (err Error) Error() string {
	return err.Wrapped.Message
}
