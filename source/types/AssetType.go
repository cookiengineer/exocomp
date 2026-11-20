package types

type AssetType string

const (
	AssetTypeIP          AssetType = "ip"
	AssetTypeCIDR        AssetType = "cidr"
	AssetTypeDomain      AssetType = "domain"
	AssetTypeSubdomain   AssetType = "subdomain"
	AssetTypeURL         AssetType = "url"
	AssetTypeEndpoint    AssetType = "endpoint"
	AssetTypePort        AssetType = "port"
	AssetTypeService     AssetType = "service"
	AssetTypeTechnology  AssetType = "technology"
	AssetTypeDNSRecord   AssetType = "dns_record"
	AssetTypeCertificate AssetType = "certificate"
	AssetTypeCloudAsset  AssetType = "cloud_asset"
	AssetTypeCredential  AssetType = "credential"
)

func (asset_type AssetType) IsValid() bool {

	switch asset_type {

	case AssetTypeIP:
		return true
	case AssetTypeCIDR:
		return true
	case AssetTypeDomain:
		return true
	case AssetTypeSubdomain:
		return true
	case AssetTypeURL:
		return true
	case AssetTypeEndpoint:
		return true
	case AssetTypePort:
		return true
	case AssetTypeService:
		return true
	case AssetTypeTechnology:
		return true
	case AssetTypeDNSRecord:
		return true
	case AssetTypeCertificate:
		return true
	case AssetTypeCloudAsset:
		return true
	case AssetTypeCredential:
		return true
	default:
		return false

	}

}
