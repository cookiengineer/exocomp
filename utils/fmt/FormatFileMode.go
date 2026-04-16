package utils

import "fmt"
import "os"
import "strings"

func FormatFileMode(mode os.FileMode) string {

	file_type   := ""
	permissions := []string{}

	if mode&0400 != 0 {
		permissions = append(permissions, "readable")
	}
	if mode&0200 != 0 {
		permissions = append(permissions, "writable")
	}
	if mode&0100 != 0 {
		permissions = append(permissions, "executable")
	}

	switch {
	case mode.IsDir():
		file_type = "directory"
	case mode&os.ModeSymlink != 0:
		file_type = "symlink"
	default:
		file_type = "file"
	}

	if len(permissions) == 0 {
		return file_type
	}

	return fmt.Sprintf("%s (%s)", file_type, strings.Join(permissions, ", "))

}

