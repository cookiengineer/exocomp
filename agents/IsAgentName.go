package agents

import "exocomp/utils"

func IsAgentName(raw string) bool {

	return raw == utils.FormatAgentName(raw)

}
