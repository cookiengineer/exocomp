package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func updateAsset(tool *Assets, id string, asset types.Asset, arguments utils_chat.Arguments) (string, error) {

	if type_value, ok := arguments.String("type"); ok == true {
		asset.Type = types.AssetType(strings.ToLower(strings.TrimSpace(type_value)))
	}

	if value, ok := arguments.String("value"); ok == true {
		asset.Value = value
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

	if err := validateAsset(&asset); err == nil {

		conflict := ""

		for other_id, other := range tool.contents {

			if other_id != id && assetDedupKey(other) == assetDedupKey(asset) {
				conflict = other_id
				break
			}

		}

		if conflict == "" {

			asset.LastSeen = schemas.NewDatetime()
			tool.contents[id] = asset

			if err1 := writeAssetsIndex(tool); err1 == nil {
				return fmt.Sprintf("assets.Update: Asset \"%s\" updated.", id), nil
			} else {
				return "", fmt.Errorf("assets.Update: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Update: Asset value conflicts with \"%s\".", conflict)
		}

	} else {
		return "", fmt.Errorf("assets.Update: %s", err.Error())
	}

}
