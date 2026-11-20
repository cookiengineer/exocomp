package types

import "exocomp/schemas"

type Asset struct {
	ID         string            `json:"id"`
	Type       AssetType         `json:"type"`
	Value      string            `json:"value"`
	ParentID   string            `json:"parent_id,omitempty"`
	Domain     string            `json:"domain,omitempty"`
	IP         string            `json:"ip,omitempty"`
	Port       int               `json:"port,omitempty"`
	URL        string            `json:"url,omitempty"`
	Endpoint   string            `json:"endpoint,omitempty"`
	Technology []string          `json:"technology,omitempty"`
	Props      map[string]string `json:"props,omitempty"`
	Source     string            `json:"source,omitempty"`
	Confidence FindingConfidence `json:"confidence,omitempty"`
	FirstSeen  schemas.Datetime  `json:"first_seen"`
	LastSeen   schemas.Datetime  `json:"last_seen"`
	InScope    bool              `json:"in_scope"`
	Tested     bool              `json:"tested"`
	Tags       []string          `json:"tags,omitempty"`
}

func (asset *Asset) Validate() bool {

	if IsAssetID(asset.ID) == false {
		return false
	} else if asset.Type.IsValid() == false {
		return false
	} else if asset.Value == "" {
		return false
	} else if asset.Port < 0 || asset.Port > 65535 {
		return false
	}

	if asset.Confidence == "" || asset.Confidence.IsValid() == true {

		for _, technology := range asset.Technology {
			if technology == "" {
				return false
			}
		}

		return true

	} else {
		return false
	}
}
