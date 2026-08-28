package tools

import "exocomp/types"
import "fmt"
import "strings"

func getWebsiteUserAgent(id string) (*types.UserAgent, error) {

	if strings.TrimSpace(id) == "" {
		id = "chrome-windows"
	}

	useragent := types.GetUserAgent(id)

	if useragent != nil {
		return useragent, nil
	} else {
		return nil, fmt.Errorf("Invalid User-Agent \"%s\"", id)
	}

}

