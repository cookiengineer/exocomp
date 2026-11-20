package types

type FindingStatus string

const (
	FindingStatusCandidate     FindingStatus = "candidate"
	FindingStatusVerified      FindingStatus = "verified"
	FindingStatusFixed         FindingStatus = "fixed"
	FindingStatusFalsePositive FindingStatus = "false_positive"
	FindingStatusAcceptedRisk  FindingStatus = "accepted_risk"
)

func (status FindingStatus) IsValid() bool {

	switch status {

	case FindingStatusCandidate:
		return true
	case FindingStatusVerified:
		return true
	case FindingStatusFixed:
		return true
	case FindingStatusFalsePositive:
		return true
	case FindingStatusAcceptedRisk:
		return true
	default:
		return false

	}

}
