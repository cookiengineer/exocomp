package types

import "strings"
import "time"

type ConsoleMessage struct {
	Time     time.Time `json:"time"`
	Method   string    `json:"method"`
	Value    string    `json:"value"`
	Caller   struct {
		File string `json:"file"`
		Line int    `json:"line"`
	} `json:"caller"`
}

func (message *ConsoleMessage) Lines() []string {

	result := make([]string, 0)

	if strings.Contains(message.Value, "\n") {

		tmp := strings.Split(message.Value, "\n")

		for _, line := range tmp {

			separator := ""

			if strings.HasPrefix(line, ">") {
				separator = "-"
			} else if strings.HasPrefix(line, "-") {
				separator = "-"
			} else {
				separator = " "
			}

			result = append(result, separator + line)

		}

	} else {

		tmp := []string{message.Value}

		for _, line := range tmp {

			separator := ""

			if strings.HasPrefix(line, ">") {
				separator = "-"
			} else if strings.HasPrefix(line, "-") {
				separator = "-"
			} else {
				separator = " "
			}

			result = append(result, separator + line)

		}

	}

	return result

}
