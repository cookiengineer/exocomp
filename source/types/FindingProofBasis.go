package types

type FindingProofBasis string

const (
	FindingProofBasisExecutedLive     FindingProofBasis = "executed_live"
	FindingProofBasisExecutedLab      FindingProofBasis = "executed_lab"
	FindingProofBasisStaticUnexecuted FindingProofBasis = "static_unexecuted"
)

func (basis FindingProofBasis) IsValid() bool {

	switch basis {
	case FindingProofBasisExecutedLive:
		return true
	case FindingProofBasisExecutedLab:
		return true
	case FindingProofBasisStaticUnexecuted:
		return true
	default:
		return false
	}

}
