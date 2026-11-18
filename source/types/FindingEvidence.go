package types

import "exocomp/schemas"

type FindingEvidence struct {
	ID          string              `json:"id"`
	Kind        FindingEvidenceKind `json:"kind"`
	Title       string              `json:"title"`
	Source      string              `json:"source,omitempty"`
	Command     string              `json:"command,omitempty"`
	ExitCode    *int                `json:"exit_code,omitempty"`
	Environment map[string]string   `json:"environment,omitempty"`
	Content     string              `json:"content,omitempty"`
	Artifact    string              `json:"artifact,omitempty"`
	SHA256      string              `json:"sha256,omitempty"`
	Size        int64               `json:"size,omitempty"`
	ContentType string              `json:"content_type,omitempty"`
	CreatedAt   schemas.Datetime    `json:"created_at"`
}
