package assets

import "exocomp/types"
import "fmt"

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
