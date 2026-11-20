package types

type FindingConfidence string

const (
	FindingConfidenceObserved   FindingConfidence = "observed"
	FindingConfidenceInferred   FindingConfidence = "inferred"
	FindingConfidenceUnverified FindingConfidence = "unverified"
)

func (confidence FindingConfidence) IsValid() bool {

	switch confidence {

	case FindingConfidenceObserved:
		return true

	case FindingConfidenceInferred:
		return true

	case FindingConfidenceUnverified:
		return true

	default:
		return false

	}

}
