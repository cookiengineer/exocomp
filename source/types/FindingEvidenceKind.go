package types

type FindingEvidenceKind string

const (
	FindingEvidenceKindCommandOutput FindingEvidenceKind = "command_output"
	FindingEvidenceKindHTTPRequest   FindingEvidenceKind = "http_request"
	FindingEvidenceKindHTTPResponse  FindingEvidenceKind = "http_response"
	FindingEvidenceKindArtifact      FindingEvidenceKind = "artifact"
	FindingEvidenceKindNote          FindingEvidenceKind = "note"
	FindingEvidenceKindUnitTest      FindingEvidenceKind = "unit_test"
	FindingEvidenceKindPoc           FindingEvidenceKind = "poc"
	FindingEvidenceKindTestRun       FindingEvidenceKind = "test_run"
)

func (kind FindingEvidenceKind) IsValid() bool {

	switch kind {
	case FindingEvidenceKindCommandOutput:
		return true
	case FindingEvidenceKindHTTPRequest:
		return true
	case FindingEvidenceKindHTTPResponse:
		return true
	case FindingEvidenceKindArtifact:
		return true
	case FindingEvidenceKindNote:
		return true
	case FindingEvidenceKindUnitTest:
		return true
	case FindingEvidenceKindPoc:
		return true
	case FindingEvidenceKindTestRun:
		return true
	default:
		return false
	}

}

func (kind FindingEvidenceKind) IsReproduction() bool {

	switch kind {
	case FindingEvidenceKindUnitTest:
		return true
	case FindingEvidenceKindPoc:
		return true
	default:
		return false
	}

}
