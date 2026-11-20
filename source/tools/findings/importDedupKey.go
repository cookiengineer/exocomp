package findings

import "exocomp/types"
import "strings"

func importDedupKey(finding types.Finding) string {
	return strings.ToLower(strings.Join([]string{finding.Title, finding.Endpoint, finding.Parameter}, "\x00"))
}
