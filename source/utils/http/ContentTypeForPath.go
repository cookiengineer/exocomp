package http

import "mime"
import "path/filepath"
import "strings"

// ContentTypeForPath returns the MIME type for a file path.
func ContentTypeForPath(path string) string {

	result := mime.TypeByExtension(filepath.Ext(path))

	if result != "" {
		return result
	} else {

		extension := strings.ToLower(filepath.Ext(path))

		if extension == ".jsonl" {
			return "application/x-ndjson"
		} else if extension == ".log" || extension == ".txt" {
			return "text/plain"
		}

		return "application/octet-stream"

	}
}
