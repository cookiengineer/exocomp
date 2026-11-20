package assets

import utils_chat "exocomp/utils/chat"

func assetExportLimit(arguments utils_chat.Arguments) int {

	limit := 1000

	if value, ok := arguments.Int("limit"); ok == true && value > 0 {
		limit = value
	}

	if limit > 10000 {
		limit = 10000
	}

	return limit

}
