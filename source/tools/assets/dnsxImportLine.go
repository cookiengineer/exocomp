package assets

type dnsxImportLine struct {
	Host       string   `json:"host"`
	A          []string `json:"a"`
	AAAA       []string `json:"aaaa"`
	CNAME      []string `json:"cname"`
	MX         []string `json:"mx"`
	NS         []string `json:"ns"`
	TXT        []string `json:"txt"`
	StatusCode string   `json:"status_code"`
}
