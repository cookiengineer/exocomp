package findings

import "exocomp/types"
import "sort"

func SortSummaries(summaries []types.FindingSummary) {

	sort.Slice(summaries, func(a int, b int) bool {

		rank1 := SeverityRank(summaries[a].Severity)
		rank2 := SeverityRank(summaries[b].Severity)

		if rank1 != rank2 {
			return rank1 < rank2
		} else if summaries[a].UpdatedAt.Equal(summaries[b].UpdatedAt.Time) == false {
			return summaries[a].UpdatedAt.After(summaries[b].UpdatedAt.Time)
		} else {
			return summaries[a].ID < summaries[b].ID
		}

	})

}
