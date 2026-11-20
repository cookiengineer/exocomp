package assets

import utils_chat "exocomp/utils/chat"

func assetLimit(arguments utils_chat.Arguments) int {

	limit := assetsDefaultLimit

	if value, ok := arguments.Int("limit"); ok == true && value > 0 {
		limit = value
	}

	if limit > assetsMaxLimit {
		limit = assetsMaxLimit
	}

	return limit

}
