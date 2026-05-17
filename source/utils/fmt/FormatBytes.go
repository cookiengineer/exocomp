package fmt

import "fmt"

func FormatBytes(bytes uint64) string {

	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div := uint64(unit)
	exp := 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)

}

