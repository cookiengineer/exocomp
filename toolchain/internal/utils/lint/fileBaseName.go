package lint

import "path"
import "strings"

func fileBaseName(file *File) string {
	return strings.TrimSuffix(path.Base(file.Rel), ".go")
}
