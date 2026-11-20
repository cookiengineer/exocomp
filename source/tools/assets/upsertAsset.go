package assets

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import "fmt"

func upsertAsset(tool *Assets, asset types.Asset) (types.Asset, bool, error) {

	asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(utils_findings.NormalizeAssetValue(asset.Type, asset.Value)), assetsMaxValue)
	asset.Domain = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Domain), assetsMaxValue)
	asset.IP = utils_findings.LimitText(utils_findings.RedactSecrets(asset.IP), assetsMaxValue)
	asset.URL = utils_findings.LimitText(utils_findings.RedactSecrets(asset.URL), assetsMaxValue)
	asset.Endpoint = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Endpoint), assetsMaxValue)
	asset.Source = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Source), assetsMaxText)
	asset.Technology = normalizeStringList(asset.Technology, assetsMaxText)
	asset.Tags = normalizeStringList(asset.Tags, assetsMaxText)
	asset.Props = normalizeProps(asset.Props)

	if err := validateAsset(&asset); err == nil {

		now := schemas.NewDatetime()
		key := assetDedupKey(asset)

		for id, existing := range tool.contents {

			if assetDedupKey(existing) != key {
				continue
			} else {

				existing.LastSeen = now

				if existing.ParentID == "" && asset.ParentID != "" {
					existing.ParentID = asset.ParentID
				}

				if existing.Domain == "" {
					existing.Domain = asset.Domain
				}

				if existing.IP == "" {
					existing.IP = asset.IP
				}

				if existing.Port == 0 {
					existing.Port = asset.Port
				}

				if existing.URL == "" {
					existing.URL = asset.URL
				}

				if existing.Endpoint == "" {
					existing.Endpoint = asset.Endpoint
				}

				if existing.Source == "" {
					existing.Source = asset.Source
				}

				if existing.Confidence == "" {
					existing.Confidence = asset.Confidence
				}

				existing.Technology = mergeStringList(existing.Technology, asset.Technology)
				existing.Tags = mergeStringList(existing.Tags, asset.Tags)
				existing.Props = mergeProps(existing.Props, asset.Props)
				existing.InScope = existing.InScope || asset.InScope
				existing.Tested = existing.Tested || asset.Tested

				tool.contents[id] = existing

				return existing, false, nil

			}
		}

		asset.ID = utils_findings.NewAssetID()

		if types.IsAssetID(asset.ID) == true {

			if asset.FirstSeen.IsZero() == true {
				asset.FirstSeen = now
			}

			asset.LastSeen = now

			if asset.Confidence == "" {
				asset.Confidence = types.FindingConfidenceObserved
			}

			if asset.Source == "" {
				asset.Source = "manual"
			}

			if asset.ParentID == "" {
				asset.ParentID = inferAssetParent(tool, asset)
			}

			if asset.Validate() == true {

				tool.contents[asset.ID] = asset

				return asset, true, nil

			} else {
				return types.Asset{}, false, fmt.Errorf("Invalid asset.")
			}
		} else {
			return types.Asset{}, false, fmt.Errorf("Cannot generate asset ID.")
		}
	} else {
		return types.Asset{}, false, err
	}
}
