package findings

import "exocomp/types"
import "sort"

func ReadAll(playground string) ([]types.Finding, error) {

	summaries, err0 := ReadIndex(playground)

	if err0 == nil {

		ids := make([]string, 0, len(summaries))

		for _, summary := range summaries {
			ids = append(ids, summary.ID)
		}

		sort.Strings(ids)

		result := make([]types.Finding, 0, len(ids))

		for _, id := range ids {

			finding, err1 := ReadFinding(playground, id)

			if err1 == nil {
				result = append(result, finding)
			} else {
				return nil, err1
			}

		}

		return result, nil

	} else {
		return nil, err0
	}

}
