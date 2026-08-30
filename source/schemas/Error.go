package schemas

type Error struct {
	err struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   any     `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

func (err Error) Error() string {
	return err.err.Message
}
