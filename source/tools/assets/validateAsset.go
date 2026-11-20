package assets

import "exocomp/types"
import "fmt"
import "net"
import "net/url"
import "strconv"
import "strings"

func validateAsset(asset *types.Asset) error {

	if asset.Type.IsValid() == false {
		return fmt.Errorf("Invalid asset type.")
	} else if strings.TrimSpace(asset.Value) == "" {
		return fmt.Errorf("Invalid asset value.")
	}

	switch asset.Type {
	case types.AssetTypeIP:

		if net.ParseIP(asset.Value) == nil {
			return fmt.Errorf("Invalid IP address \"%s\".", asset.Value)
		}

	case types.AssetTypeCIDR:

		if _, _, err := net.ParseCIDR(asset.Value); err != nil {
			return fmt.Errorf("Invalid CIDR \"%s\".", asset.Value)
		}

	case types.AssetTypeURL:

		parsed, err := url.Parse(asset.Value)

		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("Invalid URL \"%s\".", asset.Value)
		}

	case types.AssetTypeDomain, types.AssetTypeSubdomain:

		if strings.ContainsAny(asset.Value, " /\\") == true {
			return fmt.Errorf("Invalid domain \"%s\".", asset.Value)
		}

	case types.AssetTypePort:

		if asset.Port == 0 {

			if _, err := strconv.Atoi(asset.Value); err != nil && strings.Contains(asset.Value, ":") == false {
				return fmt.Errorf("Invalid port \"%s\".", asset.Value)
			}

		}

	}

	if asset.Port >= 0 && asset.Port <= 65535 {

		if asset.Confidence == "" || asset.Confidence.IsValid() == true {
			return nil
		} else {
			return fmt.Errorf("Invalid confidence \"%s\".", asset.Confidence)
		}
	} else {
		return fmt.Errorf("Invalid port number.")
	}
}
