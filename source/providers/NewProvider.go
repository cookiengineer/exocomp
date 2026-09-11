package providers

import "exocomp/types"

func NewProvider(model string) *types.Provider {

	preset, ok := Providers[model]

	if ok == true && preset != nil {

		clone := *preset

		return &clone

	}

	return nil

}
