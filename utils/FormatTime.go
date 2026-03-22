package utils

import "time"

func FormatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}
