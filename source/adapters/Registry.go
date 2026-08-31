package adapters

import "exocomp/types"

var Registry []types.Adapter

func init() {

	Registry = []types.Adapter{
		&DeepSeek{},
		&NoTransform{},
	}

}
