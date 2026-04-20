package chat

import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "strings"

func SummarizeMessages(chat []*schemas.Message, user bool, assistant bool, tool bool) string {

	result := make([]string, 0)

	if len(chat) > 0 {

		for _, message := range chat {

			if message.Role == "assistant" && assistant == true {

				content := strings.TrimSpace(utils_fmt.FormatSingleLine(message.Content))
				result = append(result, fmt.Sprintf("Assistant wrote: %s", content))

			} else if message.Role == "user" && user == true {

				content := strings.TrimSpace(utils_fmt.FormatSingleLine(message.Content))
				result = append(result, fmt.Sprintf("User wrote: %s", content))

			} else if message.Role == "tool" && tool == true {

				lines := strings.Split(message.Content, "\n")

				if len(lines) > 0 {
					content := strings.TrimSpace(lines[0])
					result = append(result, fmt.Sprintf("Tool executed: %s", content))
				}

			} else if message.Role == "system" {

				// Do Nothing

			}

		}

	}

	return strings.Join(result, "\n\n")

}
