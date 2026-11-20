package types

type FindingSeverity string

const (
	FindingSeverityCritical FindingSeverity = "critical"
	FindingSeverityHigh     FindingSeverity = "high"
	FindingSeverityMedium   FindingSeverity = "medium"
	FindingSeverityLow      FindingSeverity = "low"
	FindingSeverityInfo     FindingSeverity = "info"
)

func (severity FindingSeverity) IsValid() bool {

	switch severity {

	case FindingSeverityCritical:
		return true

	case FindingSeverityHigh:
		return true

	case FindingSeverityMedium:
		return true

	case FindingSeverityLow:
		return true

	case FindingSeverityInfo:
		return true

	default:
		return false

	}

}
