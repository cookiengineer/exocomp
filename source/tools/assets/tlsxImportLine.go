package assets

type tlsxImportLine struct {
	Host            string   `json:"host"`
	Port            int      `json:"port"`
	TLSVersion      string   `json:"tls_version"`
	Cipher          string   `json:"cipher"`
	SubjectCN       string   `json:"subject_cn"`
	SubjectAN       []string `json:"subject_an"`
	IssuerCN        string   `json:"issuer_cn"`
	IssuerOrg       string   `json:"issuer_org"`
	NotBefore       string   `json:"not_before"`
	NotAfter        string   `json:"not_after"`
	Serial          string   `json:"serial"`
	FingerprintHash string   `json:"fingerprint_hash"`
	Expired         bool     `json:"expired"`
	SelfSigned      bool     `json:"self_signed"`
	Mismatched      bool     `json:"mismatched"`
}
