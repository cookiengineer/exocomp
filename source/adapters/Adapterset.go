package adapters

import "exocomp/types"
import net_url "net/url"

func Adapterset(url *net_url.URL, model string) []types.Adapter {

	result_adapters := make([]types.Adapter, 0)

	for _, adapter := range Registry {

		if adapter.Detect(url, model) == true {
			result_adapters = append(result_adapters, adapter)
		}

	}

	return result_adapters

}
