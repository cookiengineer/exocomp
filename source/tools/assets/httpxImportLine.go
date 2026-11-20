package assets

type httpxImportLine struct {
	Input      string   `json:"input"`
	URL        string   `json:"url"`
	Host       string   `json:"host"`
	Scheme     string   `json:"scheme"`
	StatusCode int      `json:"status_code"`
	Title      string   `json:"title"`
	Tech       []string `json:"tech"`
	Webserver  string   `json:"webserver"`
	CDNName    string   `json:"cdn_name"`
	IP         string   `json:"ip"`
}
