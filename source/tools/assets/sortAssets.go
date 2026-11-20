package assets

import "exocomp/types"
import "sort"

func sortAssets(entries []types.Asset) {

	sort.Slice(entries, func(a int, b int) bool {

		if entries[a].Type == entries[b].Type {
			return entries[a].Value < entries[b].Value
		}

		return entries[a].Type < entries[b].Type

	})

}
