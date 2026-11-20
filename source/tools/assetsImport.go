package tools

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import utils_strings "exocomp/utils/strings"
import "encoding/json"
import "fmt"
import "net/url"
import "strconv"
import "strings"

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

type naabuImportLine struct {
	Host     string `json:"host"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	TLS      bool   `json:"tls"`
}

type subfinderImportLine struct {
	Host   string `json:"host"`
	Source string `json:"source"`
	Input  string `json:"input"`
}

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

type urlfinderImportLine struct {
	URL    string `json:"url"`
	Source string `json:"source"`
	Input  string `json:"input"`
}

func (tool *Assets) Import(arguments utils_chat.Arguments) (string, error) {

	path, ok := arguments.String("path")

	if ok == false || strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("assets.Import: Invalid parameter \"path\" is not a non-empty string.")
	}

	format, ok := arguments.String("format")

	if ok == false {
		return "", fmt.Errorf("assets.Import: Invalid parameter \"format\" is not a string.")
	}

	format = strings.ToLower(strings.TrimSpace(format))

	if isSupportedAssetFormat(format) == false {
		return "", fmt.Errorf("assets.Import: Unsupported format \"%s\".", format)
	}

	target, _ := arguments.String("target")
	in_scope, _ := arguments.Bool("in_scope")
	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	data, err0 := utils_findings.ReadScanArtifactFile(sandbox, path, "assets.Import")

	if err0 == nil {

		parsed, invalid, err1 := parseAssetResults(format, data, target, in_scope)

		if err1 == nil {

			created, updated, invalid, err2 := importAssetResults(tool, parsed, invalid)

			if err2 == nil {
				return fmt.Sprintf("assets.Import: Created %d assets, updated %d, skipped %d invalid results (format %s).", created, updated, invalid, format), nil
			} else {
				return "", fmt.Errorf("assets.Import: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Import: %s", err1.Error())
		}

	} else {
		return "", err0
	}

}

func importAssetResults(tool *Assets, parsed []types.Asset, invalid int) (int, int, int, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		created := 0
		updated := 0

		for _, asset := range parsed {

			_, is_new, err1 := upsertAsset(tool, asset)

			if err1 == nil {

				if is_new == true {
					created++
				} else {
					updated++
				}

			} else {
				invalid++
			}

		}

		if err2 := writeAssetsIndex(tool); err2 == nil {
			return created, updated, invalid, nil
		} else {
			return 0, 0, 0, err2
		}

	} else {
		return 0, 0, 0, err
	}

}

func isSupportedAssetFormat(format string) bool {

	switch format {
	case "httpx", "naabu", "subfinder", "dnsx", "tlsx", "katana", "gobuster", "urlfinder":
		return true
	default:
		return false
	}

}

func parseAssetResults(format string, data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	switch format {
	case "httpx":
		return parseHttpxResults(data, target, in_scope)
	case "naabu":
		return parseNaabuResults(data, in_scope)
	case "subfinder":
		return parseSubfinderResults(data, target, in_scope)
	case "dnsx":
		return parseDnsxResults(data, in_scope)
	case "tlsx":
		return parseTlsxResults(data, in_scope)
	case "katana":
		return parseKatanaResults(data, in_scope)
	case "gobuster":
		return parseGobusterResults(data, target, in_scope)
	case "urlfinder":
		return parseUrlfinderResults(data, in_scope)
	default:
		return nil, 0, fmt.Errorf("Unsupported format \"%s\".", format)
	}

}

func parseHttpxResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := httpxImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		value := utils_strings.FirstNonEmpty(parsed.URL, parsed.Host, parsed.Input)

		if value == "" {
			invalid++
			continue
		}

		asset_type := types.AssetTypeSubdomain

		if strings.Contains(value, "://") == true {
			asset_type = types.AssetTypeURL
		} else if target != "" && strings.EqualFold(value, target) == true {
			asset_type = types.AssetTypeDomain
		}

		asset := newImportedAsset(asset_type, value, "httpx", in_scope)
		asset.URL = parsed.URL
		asset.Endpoint = parsed.URL
		asset.IP = parsed.IP
		asset.Domain = hostFromValue(value)
		asset.Technology = append([]string{}, parsed.Tech...)
		asset.Props = map[string]string{}

		if parsed.StatusCode > 0 {
			asset.Props["status_code"] = strconv.Itoa(parsed.StatusCode)
		}

		if parsed.Title != "" {
			asset.Props["title"] = parsed.Title
		}

		if parsed.Webserver != "" {
			asset.Props["webserver"] = parsed.Webserver
		}

		if parsed.CDNName != "" {
			asset.Props["cdn"] = parsed.CDNName
		}

		if parsed.Scheme != "" {
			asset.Props["scheme"] = parsed.Scheme
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseNaabuResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := naabuImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		host := utils_strings.FirstNonEmpty(parsed.Host, parsed.IP)

		if host == "" || parsed.Port < 1 || parsed.Port > 65535 {
			invalid++
			continue
		}

		asset := newImportedAsset(types.AssetTypePort, fmt.Sprintf("%s:%d", host, parsed.Port), "naabu", in_scope)
		asset.Domain = host
		asset.IP = parsed.IP
		asset.Port = parsed.Port
		asset.Props = map[string]string{}

		if parsed.Protocol != "" {
			asset.Props["protocol"] = parsed.Protocol
		}

		if parsed.TLS == true {
			asset.Props["tls"] = "true"
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseSubfinderResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := subfinderImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Host == "" {
			invalid++
			continue
		}

		asset_type := types.AssetTypeSubdomain

		if target != "" && strings.EqualFold(parsed.Host, target) == true {
			asset_type = types.AssetTypeDomain
		}

		asset := newImportedAsset(asset_type, parsed.Host, "subfinder", in_scope)
		asset.Domain = parsed.Host
		asset.Props = map[string]string{}

		if parsed.Source != "" {
			asset.Props["source"] = parsed.Source
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseDnsxResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := dnsxImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Host == "" {
			invalid++
			continue
		}

		asset := newImportedAsset(types.AssetTypeDNSRecord, parsed.Host, "dnsx", in_scope)
		asset.Domain = parsed.Host
		asset.Props = map[string]string{}

		if len(parsed.A) > 0 {
			asset.Props["a"] = strings.Join(parsed.A, ", ")
		}

		if len(parsed.AAAA) > 0 {
			asset.Props["aaaa"] = strings.Join(parsed.AAAA, ", ")
		}

		if len(parsed.CNAME) > 0 {
			asset.Props["cname"] = strings.Join(parsed.CNAME, ", ")
		}

		if len(parsed.MX) > 0 {
			asset.Props["mx"] = strings.Join(parsed.MX, ", ")
		}

		if len(parsed.NS) > 0 {
			asset.Props["ns"] = strings.Join(parsed.NS, ", ")
		}

		if len(parsed.TXT) > 0 {
			asset.Props["txt"] = strings.Join(parsed.TXT, ", ")
		}

		if parsed.StatusCode != "" {
			asset.Props["status_code"] = parsed.StatusCode
		}

		result = append(result, asset)

		for _, ip := range parsed.A {

			ip_asset := newImportedAsset(types.AssetTypeIP, ip, "dnsx", in_scope)
			ip_asset.Domain = parsed.Host
			result = append(result, ip_asset)

		}

		for _, ip := range parsed.AAAA {

			ip_asset := newImportedAsset(types.AssetTypeIP, ip, "dnsx", in_scope)
			ip_asset.Domain = parsed.Host
			result = append(result, ip_asset)

		}

	}

	return result, invalid, nil

}

func parseTlsxResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := tlsxImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Host == "" {
			invalid++
			continue
		}

		value := utils_strings.FirstNonEmpty(parsed.Serial, parsed.FingerprintHash)

		if value == "" {
			value = fmt.Sprintf("%s:%d", parsed.Host, parsed.Port)
		}

		asset := newImportedAsset(types.AssetTypeCertificate, value, "tlsx", in_scope)
		asset.Domain = parsed.Host
		asset.Port = parsed.Port
		asset.Props = map[string]string{}

		if parsed.TLSVersion != "" {
			asset.Props["tls_version"] = parsed.TLSVersion
		}

		if parsed.Cipher != "" {
			asset.Props["cipher"] = parsed.Cipher
		}

		if parsed.SubjectCN != "" {
			asset.Props["subject_cn"] = parsed.SubjectCN
		}

		if len(parsed.SubjectAN) > 0 {
			asset.Props["subject_an"] = strings.Join(parsed.SubjectAN, ", ")
		}

		if parsed.IssuerCN != "" {
			asset.Props["issuer_cn"] = parsed.IssuerCN
		}

		if parsed.IssuerOrg != "" {
			asset.Props["issuer_org"] = parsed.IssuerOrg
		}

		if parsed.NotBefore != "" {
			asset.Props["not_before"] = parsed.NotBefore
		}

		if parsed.NotAfter != "" {
			asset.Props["not_after"] = parsed.NotAfter
		}

		if parsed.Expired == true {
			asset.Props["expired"] = "true"
		}

		if parsed.SelfSigned == true {
			asset.Props["self_signed"] = "true"
		}

		if parsed.Mismatched == true {
			asset.Props["mismatched"] = "true"
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseKatanaResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := katanaImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Request.Endpoint == "" {
			invalid++
			continue
		}

		asset := newImportedAsset(types.AssetTypeEndpoint, parsed.Request.Endpoint, "katana", in_scope)
		asset.URL = parsed.Request.Endpoint
		asset.Endpoint = parsed.Request.Endpoint
		asset.Domain = hostFromValue(parsed.Request.Endpoint)
		asset.Props = map[string]string{}

		if parsed.Request.Method != "" {
			asset.Props["method"] = parsed.Request.Method
		}

		if parsed.Response.StatusCode > 0 {
			asset.Props["status_code"] = strconv.Itoa(parsed.Response.StatusCode)
		}

		if parsed.Request.Tag != "" {
			asset.Props["tag"] = parsed.Request.Tag
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseGobusterResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	result := make([]types.Asset, 0, len(lines))
	invalid := 0

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		marker := strings.Index(line, " (Status:")

		if marker < 0 {
			continue
		}

		value := strings.TrimSpace(line[:marker])
		status := ""

		if relative_end := strings.Index(line[marker:], ")"); relative_end >= 0 {
			status = strings.TrimSpace(line[marker+len(" (Status:") : marker+relative_end])
		}

		if value == "" {
			invalid++
			continue
		}

		if strings.HasPrefix(value, "http://") == false && strings.HasPrefix(value, "https://") == false {

			if target != "" {
				value = strings.TrimRight(target, "/") + "/" + strings.TrimLeft(value, "/")
			}

		}

		asset := newImportedAsset(types.AssetTypeEndpoint, value, "gobuster", in_scope)
		asset.URL = value
		asset.Endpoint = value
		asset.Domain = hostFromValue(value)
		asset.Props = map[string]string{}

		if status != "" {
			asset.Props["status_code"] = status
		}

		if size_marker := strings.Index(line, "[Size:"); size_marker >= 0 {

			if end := strings.Index(line[size_marker:], "]"); end >= 0 {
				size := strings.TrimSpace(line[size_marker+len("[Size:") : size_marker+end])
				asset.Props["size"] = size
			}

		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func parseUrlfinderResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := urlfinderImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.URL == "" {
			invalid++
			continue
		}

		asset := newImportedAsset(types.AssetTypeURL, parsed.URL, "urlfinder", in_scope)
		asset.URL = parsed.URL
		asset.Endpoint = parsed.URL
		asset.Domain = hostFromValue(parsed.URL)
		asset.Props = map[string]string{}

		if parsed.Source != "" {
			asset.Props["source"] = parsed.Source
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}

func decodeImportLines(data []byte) ([]json.RawMessage, int, error) {

	trimmed := strings.TrimSpace(string(data))

	if trimmed == "" {
		return make([]json.RawMessage, 0), 0, nil
	}

	if strings.HasPrefix(trimmed, "[") == true {

		result := make([]json.RawMessage, 0)

		if err := json.Unmarshal([]byte(trimmed), &result); err != nil {
			return nil, 0, fmt.Errorf("Cannot parse results.")
		}

		return result, 0, nil

	}

	result := make([]json.RawMessage, 0)
	invalid := 0

	for _, line := range strings.Split(trimmed, "\n") {

		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") == true {
			continue
		}

		if json.Valid([]byte(line)) == false {
			invalid++
			continue
		}

		result = append(result, json.RawMessage(line))

	}

	return result, invalid, nil

}

func newImportedAsset(asset_type types.AssetType, value string, source string, in_scope bool) types.Asset {

	return types.Asset{
		Type:       asset_type,
		Value:      utils_findings.LimitText(value, assetsMaxValue),
		Source:     source,
		Confidence: types.FindingConfidenceObserved,
		InScope:    in_scope,
	}
}

func hostFromValue(value string) string {

	if strings.Contains(value, "://") == true {

		if parsed, err := url.Parse(value); err == nil {
			return parsed.Hostname()
		}

	}

	if strings.Contains(value, "/") == true {
		return strings.Split(value, "/")[0]
	}

	return value

}
