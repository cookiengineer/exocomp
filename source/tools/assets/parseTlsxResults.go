package assets

import "exocomp/types"
import utils_strings "exocomp/utils/strings"
import "encoding/json"
import "fmt"
import "strings"

func parseTlsxResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err == nil {

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

	} else {
		return nil, 0, err
	}

}
