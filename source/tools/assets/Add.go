package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Assets) Add(arguments utils_chat.Arguments) (string, error) {

	type_value, has_type := arguments.String("type")
	value, has_value := arguments.String("value")

	if has_type == false || has_value == false {
		return "", fmt.Errorf("assets.Add: Invalid parameters \"type\" and \"value\" are required.")
	}

	asset := types.Asset{
		Type:       types.AssetType(strings.ToLower(strings.TrimSpace(type_value))),
		Value:      value,
		Technology: make([]string, 0),
		Props:      make(map[string]string),
		Tags:       make([]string, 0),
	}

	if asset.Type.IsValid() == false {
		return "", fmt.Errorf("assets.Add: Invalid asset type \"%s\".", type_value)
	}

	if parent_id, ok := arguments.String("parent_id"); ok == true {
		asset.ParentID = parent_id
	}

	if domain, ok := arguments.String("domain"); ok == true {
		asset.Domain = utils_findings.LimitText(utils_findings.RedactSecrets(domain), assetsMaxValue)
	}

	if ip, ok := arguments.String("ip"); ok == true {
		asset.IP = utils_findings.LimitText(utils_findings.RedactSecrets(ip), assetsMaxValue)
	}

	if port, ok := arguments.Int("port"); ok == true {
		asset.Port = port
	}

	if raw_url, ok := arguments.String("url"); ok == true {
		asset.URL = utils_findings.LimitText(utils_findings.RedactSecrets(raw_url), assetsMaxValue)
	}

	if endpoint, ok := arguments.String("endpoint"); ok == true {
		asset.Endpoint = utils_findings.LimitText(utils_findings.RedactSecrets(endpoint), assetsMaxValue)
	}

	if technology, ok := arguments.Strings("technology"); ok == true {
		asset.Technology = normalizeStringList(technology, assetsMaxText)
	}

	if source, ok := arguments.String("source"); ok == true {
		asset.Source = utils_findings.LimitText(utils_findings.RedactSecrets(source), assetsMaxText)
	}

	if confidence, ok := arguments.String("confidence"); ok == true {
		asset.Confidence = types.FindingConfidence(strings.ToLower(strings.TrimSpace(confidence)))
	}

	if in_scope, ok := arguments.Bool("in_scope"); ok == true {
		asset.InScope = in_scope
	}

	if tested, ok := arguments.Bool("tested"); ok == true {
		asset.Tested = tested
	}

	if tags, ok := arguments.Strings("tags"); ok == true {
		asset.Tags = normalizeStringList(tags, assetsMaxText)
	}

	if props, ok := arguments.Map("props"); ok == true {
		asset.Props = normalizeProps(props)
	}

	asset.Value = utils_findings.NormalizeAssetValue(asset.Type, asset.Value)
	asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Value), assetsMaxValue)

	if err := validateAsset(&asset); err != nil {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

	stored, created, err := storeAsset(tool, asset)

	if err == nil {

		if created == true {
			return fmt.Sprintf("assets.Add: Asset \"%s\" created.", stored.ID), nil
		} else {
			return fmt.Sprintf("assets.Add: Asset \"%s\" updated.", stored.ID), nil
		}

	} else {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

}
