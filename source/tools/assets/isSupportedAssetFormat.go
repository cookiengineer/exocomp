package assets

func isSupportedAssetFormat(format string) bool {

	switch format {

	case "httpx", "naabu", "subfinder", "dnsx", "tlsx", "katana", "gobuster", "urlfinder":
		return true
	default:
		return false

	}

}
