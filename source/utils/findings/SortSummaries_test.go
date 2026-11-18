package findings

import "exocomp/schemas"
import "exocomp/types"
import "testing"
import "time"

func TestSortSummaries(t *testing.T) {

	now := time.Now()

	summaries := []types.FindingSummary{
		{
			ID:        "fnd_00000001",
			Severity:  types.FindingSeverityLow,
			UpdatedAt: schemas.Datetime{Time: now.Add(2 * time.Hour)},
		},
		{
			ID:        "fnd_00000002",
			Severity:  types.FindingSeverityCritical,
			UpdatedAt: schemas.Datetime{Time: now},
		},
		{
			ID:        "fnd_00000003",
			Severity:  types.FindingSeverityHigh,
			UpdatedAt: schemas.Datetime{Time: now.Add(1 * time.Hour)},
		},
	}

	SortSummaries(summaries)

	if summaries[0].Severity != types.FindingSeverityCritical {
		t.Errorf("Expected critical finding first, got %s", summaries[0].Severity)
	} else if summaries[1].Severity != types.FindingSeverityHigh {
		t.Errorf("Expected high finding second, got %s", summaries[1].Severity)
	} else if summaries[2].Severity != types.FindingSeverityLow {
		t.Errorf("Expected low finding third, got %s", summaries[2].Severity)
	}

}
