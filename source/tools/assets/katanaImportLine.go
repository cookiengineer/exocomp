package assets

type katanaImportLine struct {
	Request struct {
		Method   string `json:"method"`
		Endpoint string `json:"endpoint"`
		Tag      string `json:"tag"`
	} `json:"request"`
	Response struct {
		StatusCode int `json:"status_code"`
	} `json:"response"`
}
