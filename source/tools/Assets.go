package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_assets "exocomp/utils/assets"
import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "net"
import "net/url"
import "slices"
import "sort"
import "strconv"
import "strings"
import "sync"

const (
	assetsDefaultLimit = 50
	assetsMaxLimit     = 200
	assetsMaxValue     = 2048
	assetsMaxText      = 1024
)

type assetArguments = findingArguments

type Assets struct {
	Methods    []string
	Playground string
	Sandbox    string
	mutex      *sync.RWMutex
	contents   map[string]types.Asset
}

func NewAssets(methods []string, playground string, sandbox string) *Assets {

	tool := &Assets{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]types.Asset),
	}

	readAssetsIndex(tool)

	return tool

}

func (tool *Assets) Name() string {
	return "assets"
}

func (tool *Assets) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		parsed := assetArguments(arguments)

		if method == "Add" {

			return tool.Add(parsed)

		} else if method == "Export" {

			return tool.Export(parsed)

		} else if method == "Import" {

			return tool.Import(parsed)

		} else if method == "Get" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Get(id)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "List" {

			return tool.List(parsed)

		} else if method == "Search" {

			return tool.Search(parsed)

		} else if method == "Update" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Update(id, parsed)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Neighbors" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Neighbors(id, parsed)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Scope" {

			return tool.Scope(parsed)

		} else if method == "Coverage" {

			return tool.Coverage()

		} else {
			return "", fmt.Errorf("assets.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("assets.%s: Method not allowed.", method)
	}

}

func (tool *Assets) Add(arguments assetArguments) (string, error) {

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

	asset.Value = utils_assets.NormalizeValue(asset.Type, asset.Value)
	asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Value), assetsMaxValue)

	if err := validateAsset(&asset); err != nil {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

	stored, created, err := upsertAsset(tool, asset)

	if err != nil {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

	if err := writeAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Add: %s", err.Error())
	}

	if created == true {
		return fmt.Sprintf("assets.Add: Asset \"%s\" created.", stored.ID), nil
	} else {
		return fmt.Sprintf("assets.Add: Asset \"%s\" updated.", stored.ID), nil
	}

}

func (tool *Assets) Get(id string) (string, error) {

	if types.IsAssetID(id) == false {
		return "", fmt.Errorf("assets.Get: Invalid asset ID \"%s\".", id)
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Get: %s", err.Error())
	}

	asset, ok := tool.contents[id]

	if ok == false {
		return "", fmt.Errorf("assets.Get: Asset \"%s\" does not exist.", id)
	}

	buffer, err := json.MarshalIndent(asset, "", "\t")

	if err != nil {
		return "", fmt.Errorf("assets.Get: Cannot serialize asset \"%s\".", id)
	}

	return strings.Join([]string{fmt.Sprintf("assets.Get: %s", id), string(buffer)}, "\n"), nil

}

func (tool *Assets) List(arguments assetArguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	query, has_query := arguments.String("query")
	in_scope_filter, has_in_scope := arguments.Bool("in_scope")
	tested_filter, has_tested := arguments.Bool("tested")
	limit := assetLimit(arguments)

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.List: Invalid parameter \"type\".")
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.List: %s", err.Error())
	}

	matches := make([]types.Asset, 0)

	for _, asset := range tool.contents {

		if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
			continue
		} else if has_query == true && assetMatchesQuery(asset, query) == false {
			continue
		} else if has_in_scope == true && asset.InScope != in_scope_filter {
			continue
		} else if has_tested == true && asset.Tested != tested_filter {
			continue
		}

		matches = append(matches, asset)

	}

	sortAssets(matches)

	total := len(matches)

	if len(matches) > limit {
		matches = matches[0:limit]
	}

	lines := make([]string, 0, len(matches))

	for _, asset := range matches {
		lines = append(lines, formatAssetSummary(asset))
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("assets.List: %d assets (showing %d).", total, len(lines)))
	result = append(result, lines...)

	return strings.Join(result, "\n"), nil

}

func (tool *Assets) Search(arguments assetArguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == false || strings.TrimSpace(query) == "" {
		return "", fmt.Errorf("assets.Search: Invalid parameter \"query\" is not a non-empty string.")
	}

	type_filter, has_type := arguments.String("type")
	limit := assetLimit(arguments)

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.Search: Invalid parameter \"type\".")
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Search: %s", err.Error())
	}

	matches := make([]types.Asset, 0)

	for _, asset := range tool.contents {

		if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
			continue
		} else if assetMatchesQuery(asset, query) == false {
			continue
		}

		matches = append(matches, asset)

	}

	sortAssets(matches)

	total := len(matches)

	if len(matches) > limit {
		matches = matches[0:limit]
	}

	lines := make([]string, 0, len(matches))

	for _, asset := range matches {
		lines = append(lines, formatAssetSummary(asset))
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("assets.Search: %d matches (showing %d).", total, len(lines)))
	result = append(result, lines...)

	return strings.Join(result, "\n"), nil

}

func (tool *Assets) Update(id string, arguments assetArguments) (string, error) {

	if types.IsAssetID(id) == false {
		return "", fmt.Errorf("assets.Update: Invalid asset ID \"%s\".", id)
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Update: %s", err.Error())
	}

	asset, ok := tool.contents[id]

	if ok == false {
		return "", fmt.Errorf("assets.Update: Asset \"%s\" does not exist.", id)
	}

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

	asset.Value = utils_assets.NormalizeValue(asset.Type, asset.Value)
	asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Value), assetsMaxValue)

	if err := validateAsset(&asset); err != nil {
		return "", fmt.Errorf("assets.Update: %s", err.Error())
	}

	for other_id, other := range tool.contents {

		if other_id != id && assetDedupKey(other) == assetDedupKey(asset) {
			return "", fmt.Errorf("assets.Update: Asset value conflicts with \"%s\".", other_id)
		}

	}

	asset.LastSeen = schemas.NewDatetime()
	tool.contents[id] = asset

	if err := writeAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Update: %s", err.Error())
	}

	return fmt.Sprintf("assets.Update: Asset \"%s\" updated.", id), nil

}

func (tool *Assets) Neighbors(id string, arguments assetArguments) (string, error) {

	if types.IsAssetID(id) == false {
		return "", fmt.Errorf("assets.Neighbors: Invalid asset ID \"%s\".", id)
	}

	direction := "both"

	if value, ok := arguments.String("direction"); ok == true {
		direction = strings.ToLower(strings.TrimSpace(value))
	}

	if direction != "parent" && direction != "children" && direction != "both" {
		return "", fmt.Errorf("assets.Neighbors: Invalid parameter \"direction\".")
	}

	limit := assetLimit(arguments)

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Neighbors: %s", err.Error())
	}

	asset, ok := tool.contents[id]

	if ok == false {
		return "", fmt.Errorf("assets.Neighbors: Asset \"%s\" does not exist.", id)
	}

	lines := make([]string, 0)

	if (direction == "parent" || direction == "both") && asset.ParentID != "" {

		if parent, ok := tool.contents[asset.ParentID]; ok == true {
			lines = append(lines, "Parent: "+formatAssetSummary(parent))
		}

	}

	if direction == "children" || direction == "both" {

		children := make([]types.Asset, 0)

		for _, candidate := range tool.contents {

			if candidate.ParentID == id {
				children = append(children, candidate)
			}

		}

		sortAssets(children)

		if len(children) > limit {
			children = children[0:limit]
		}

		if len(children) > 0 {
			lines = append(lines, fmt.Sprintf("Children: %d", len(children)))
		}

		for _, child := range children {
			lines = append(lines, formatAssetSummary(child))
		}

	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("assets.Neighbors: %s", id))

	if len(lines) == 0 {
		result = append(result, "No neighbors.")
	} else {
		result = append(result, lines...)
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Assets) Scope(arguments assetArguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	value_filter, has_value := arguments.String("value")
	query, has_query := arguments.String("query")
	has_filter := has_type == true || has_value == true || has_query == true
	in_scope := true

	if provided, ok := arguments.Bool("in_scope"); ok == true {
		in_scope = provided
	}

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.Scope: Invalid parameter \"type\".")
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Scope: %s", err.Error())
	}

	if has_filter == false {

		total := len(tool.contents)
		scoped := 0
		tested := 0

		for _, asset := range tool.contents {

			if asset.InScope == true {
				scoped++
			}

			if asset.Tested == true {
				tested++
			}

		}

		return fmt.Sprintf("assets.Scope: %d assets, %d in scope, %d tested.", total, scoped, tested), nil

	}

	updated := 0

	for id, asset := range tool.contents {

		if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
			continue
		} else if has_value == true && strings.EqualFold(asset.Value, strings.TrimSpace(value_filter)) == false {
			continue
		} else if has_query == true && assetMatchesQuery(asset, query) == false {
			continue
		}

		asset.InScope = in_scope
		asset.LastSeen = schemas.NewDatetime()
		tool.contents[id] = asset
		updated++

	}

	if err := writeAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Scope: %s", err.Error())
	}

	return fmt.Sprintf("assets.Scope: %d assets marked in_scope=%t.", updated, in_scope), nil

}

func (tool *Assets) Coverage() (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return "", fmt.Errorf("assets.Coverage: %s", err.Error())
	}

	counts := make(map[types.AssetType][3]int)

	for _, asset := range tool.contents {

		count := counts[asset.Type]

		if asset.InScope == true {
			count[0]++
		}

		if asset.Tested == true {
			count[1]++
		}

		count[2]++
		counts[asset.Type] = count

	}

	types_list := make([]string, 0, len(counts))

	for asset_type := range counts {
		types_list = append(types_list, string(asset_type))
	}

	sort.Strings(types_list)

	lines := make([]string, 0, len(types_list))

	for _, name := range types_list {

		count := counts[types.AssetType(name)]
		lines = append(lines, fmt.Sprintf("%s: %d assets, %d in scope, %d tested", name, count[2], count[0], count[1]))

	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("assets.Coverage: %d assets across %d types.", len(tool.contents), len(types_list)))
	result = append(result, lines...)

	return strings.Join(result, "\n"), nil

}

func (tool *Assets) GetContent(id string) (any, error) {

	if types.IsAssetID(id) == false {
		return nil, fmt.Errorf("assets.GetContent: Invalid asset ID \"%s\".", id)
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return nil, fmt.Errorf("assets.GetContent: %s", err.Error())
	}

	asset, ok := tool.contents[id]

	if ok == false {
		return nil, fmt.Errorf("assets.GetContent: Asset \"%s\" does not exist.", id)
	}

	return asset, nil

}

func (tool *Assets) GetContentIdentifiers() []string {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	readAssetsIndex(tool)

	result := make([]string, 0, len(tool.contents))

	for id := range tool.contents {
		result = append(result, id)
	}

	sort.Strings(result)

	return result

}

func (tool *Assets) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Assets) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err != nil {
		return nil, fmt.Errorf("assets.MarshalJSON: %s", err.Error())
	}

	entries := make([]types.Asset, 0, len(tool.contents))

	for _, asset := range tool.contents {
		entries = append(entries, asset)
	}

	sortAssets(entries)

	return json.MarshalIndent(entries, "", "\t")

}

func (tool *Assets) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range AssetsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}

func readAssetsIndex(tool *Assets) error {

	entries, err := utils_assets.ReadIndex(tool.Playground)

	if err != nil {
		return err
	}

	for id := range tool.contents {
		delete(tool.contents, id)
	}

	for _, asset := range entries {
		tool.contents[asset.ID] = asset
	}

	return nil

}

func writeAssetsIndex(tool *Assets) error {

	entries := make([]types.Asset, 0, len(tool.contents))

	for _, asset := range tool.contents {
		entries = append(entries, asset)
	}

	return utils_assets.WriteIndex(tool.Playground, entries)

}

func upsertAsset(tool *Assets, asset types.Asset) (types.Asset, bool, error) {

	asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(utils_assets.NormalizeValue(asset.Type, asset.Value)), assetsMaxValue)
	asset.Domain = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Domain), assetsMaxValue)
	asset.IP = utils_findings.LimitText(utils_findings.RedactSecrets(asset.IP), assetsMaxValue)
	asset.URL = utils_findings.LimitText(utils_findings.RedactSecrets(asset.URL), assetsMaxValue)
	asset.Endpoint = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Endpoint), assetsMaxValue)
	asset.Source = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Source), assetsMaxText)
	asset.Technology = normalizeStringList(asset.Technology, assetsMaxText)
	asset.Tags = normalizeStringList(asset.Tags, assetsMaxText)
	asset.Props = normalizeProps(asset.Props)

	if err := validateAsset(&asset); err != nil {
		return types.Asset{}, false, err
	}

	now := schemas.NewDatetime()
	key := assetDedupKey(asset)

	for id, existing := range tool.contents {

		if assetDedupKey(existing) != key {
			continue
		}

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

	asset.ID = utils_assets.NewAssetID()

	if types.IsAssetID(asset.ID) == false {
		return types.Asset{}, false, fmt.Errorf("Cannot generate asset ID.")
	}

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

	if asset.Validate() == false {
		return types.Asset{}, false, fmt.Errorf("Invalid asset.")
	}

	tool.contents[asset.ID] = asset

	return asset, true, nil

}

func inferAssetParent(tool *Assets, asset types.Asset) string {

	if asset.Domain != "" {

		domain := strings.ToLower(strings.TrimSpace(asset.Domain))

		for _, candidate := range tool.contents {

			if candidate.Type != types.AssetTypeDomain && candidate.Type != types.AssetTypeSubdomain && candidate.Type != types.AssetTypeDNSRecord {
				continue
			}

			if strings.EqualFold(candidate.Value, domain) == true {
				return candidate.ID
			}

		}

		parts := strings.Split(domain, ".")

		if len(parts) > 2 {

			parent := strings.Join(parts[1:], ".")

			for _, candidate := range tool.contents {

				if candidate.Type != types.AssetTypeDomain && candidate.Type != types.AssetTypeSubdomain && candidate.Type != types.AssetTypeDNSRecord {
					continue
				}

				if strings.EqualFold(candidate.Value, parent) == true {
					return candidate.ID
				}

			}

		}

	}

	if asset.IP != "" {

		for _, candidate := range tool.contents {

			if candidate.Type == types.AssetTypeIP && strings.EqualFold(candidate.Value, asset.IP) == true {
				return candidate.ID
			}

		}

	}

	return ""

}

func assetDedupKey(asset types.Asset) string {
	return string(asset.Type) + "\x00" + utils_assets.NormalizeValue(asset.Type, asset.Value)
}

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

	if asset.Port < 0 || asset.Port > 65535 {
		return fmt.Errorf("Invalid port number.")
	}

	if asset.Confidence != "" && asset.Confidence.IsValid() == false {
		return fmt.Errorf("Invalid confidence \"%s\".", asset.Confidence)
	}

	return nil

}

func assetLimit(arguments assetArguments) int {

	limit := assetsDefaultLimit

	if value, ok := arguments.Int("limit"); ok == true && value > 0 {
		limit = value
	}

	if limit > assetsMaxLimit {
		limit = assetsMaxLimit
	}

	return limit

}

func assetMatchesQuery(asset types.Asset, query string) bool {

	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return true
	}

	values := []string{
		string(asset.Type),
		asset.Value,
		asset.Domain,
		asset.IP,
		asset.URL,
		asset.Endpoint,
		asset.Source,
	}

	values = append(values, asset.Technology...)
	values = append(values, asset.Tags...)

	for _, value := range asset.Props {
		values = append(values, value)
	}

	for _, value := range values {

		if strings.Contains(strings.ToLower(value), query) == true {
			return true
		}

	}

	return false

}

func formatAssetSummary(asset types.Asset) string {

	flags := make([]string, 0)

	if asset.InScope == true {
		flags = append(flags, "in_scope")
	}

	if asset.Tested == true {
		flags = append(flags, "tested")
	}

	suffix := ""

	if len(flags) > 0 {
		suffix = " (" + strings.Join(flags, ", ") + ")"
	}

	return fmt.Sprintf("%s [%s] %s%s", asset.ID, asset.Type, asset.Value, suffix)

}

func sortAssets(entries []types.Asset) {

	sort.Slice(entries, func(a int, b int) bool {

		if entries[a].Type == entries[b].Type {
			return entries[a].Value < entries[b].Value
		}

		return entries[a].Type < entries[b].Type

	})

}

func normalizeStringList(values []string, maximum int) []string {

	result := make([]string, 0, len(values))

	for _, value := range values {

		value = strings.TrimSpace(utils_findings.RedactSecrets(value))

		if value == "" || slices.Contains(result, value) == true {
			continue
		}

		result = append(result, utils_findings.LimitText(value, maximum))

	}

	sort.Strings(result)

	return result

}

func normalizeProps(props map[string]string) map[string]string {

	result := make(map[string]string, len(props))

	for key, value := range props {

		key = strings.TrimSpace(utils_findings.RedactSecrets(key))
		value = strings.TrimSpace(utils_findings.RedactSecrets(value))

		if key == "" || value == "" {
			continue
		}

		result[utils_findings.LimitText(key, assetsMaxText)] = utils_findings.LimitText(value, assetsMaxText)

	}

	return result

}

func mergeStringList(left []string, right []string) []string {

	result := append([]string{}, left...)

	for _, value := range right {

		if value == "" || slices.Contains(result, value) == true {
			continue
		}

		result = append(result, value)

	}

	sort.Strings(result)

	return result

}

func mergeProps(left map[string]string, right map[string]string) map[string]string {

	result := make(map[string]string, len(left)+len(right))

	for key, value := range left {
		result[key] = value
	}

	for key, value := range right {

		if value == "" {
			continue
		}

		result[key] = value

	}

	return result

}
