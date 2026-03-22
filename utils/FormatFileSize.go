package utils

import "fmt"

func FormatFileSize(bytes int64) string {

	unit := int64(1000)

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := "kMGTPE"

	return fmt.Sprintf("%.1f %cB", float64(bytes) / float64(div), units[exp])

}

