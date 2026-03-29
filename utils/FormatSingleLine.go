package utils

import "strings"

func FormatSingleLine(text string) string {

	result := strings.TrimSpace(text)

	result = strings.ReplaceAll(result, "\r", "")
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\t", " ")

	return strings.TrimSpace(result)

}
