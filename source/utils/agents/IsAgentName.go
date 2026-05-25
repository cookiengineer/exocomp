package agents

import utils_fmt "exocomp/utils/fmt"

func IsAgentName(raw string) bool {

	return raw == utils_fmt.FormatAgentName(raw)

}
