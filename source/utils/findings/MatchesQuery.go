package findings

import "exocomp/types"
import "strings"

func MatchesQuery(finding types.Finding, query string) bool {

	values := []string{
		finding.ID,
		finding.Title,
		finding.CWE,
		finding.OWASP,
		finding.Endpoint,
		finding.Parameter,
		finding.Preconditions,
		finding.Reproduction,
		finding.Impact,
		finding.Remediation,
		finding.Counterevidence,
	}

	for _, evidence := range finding.Evidence {
		values = append(values, evidence.Title, evidence.Source, evidence.Command, evidence.Content, evidence.Artifact)
	}

	for _, value := range values {

		if strings.Contains(strings.ToLower(value), query) == true {
			return true
		}

	}

	return false

}
