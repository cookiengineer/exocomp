package agent

import utils_fmt "exocomp/utils/fmt"

func IsName(raw string) bool {

	return raw == utils_fmt.FormatAgentName(raw)

}
