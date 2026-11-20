package findings

type nucleiImportClassification struct {
	CVEID     []string `json:"cve-id"`
	CWEID     []string `json:"cwe-id"`
	CVSSScore float64  `json:"cvss-score"`
}
