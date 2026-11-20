package findings

import "exocomp/types"

func SeverityRank(severity types.FindingSeverity) int {

	switch severity {

	case types.FindingSeverityCritical:
		return 0

	case types.FindingSeverityHigh:
		return 1

	case types.FindingSeverityMedium:
		return 2

	case types.FindingSeverityLow:
		return 3

	case types.FindingSeverityInfo:
		return 4

	default:
		return 5

	}

}
