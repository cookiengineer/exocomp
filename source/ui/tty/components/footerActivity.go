package components

import "strings"

func footerActivity(status string) (string, bool) {

	trimmed := sanitizeFooterStatus(status)

	if trimmed != "" {

		lower := strings.ToLower(trimmed)

		if strings.Contains(lower, "waiting for the llm") || strings.Contains(lower, "session is busy") {
			return "Agent is working…", true
		} else if strings.Contains(lower, "answering human question") == true {
			return "Sending your answer…", true
		} else if strings.HasPrefix(lower, "running ") == true {
			return trimmed, true
		} else {
			return trimmed, false
		}

	} else {
		return "Waiting for your instructions", false
	}

}
