package markdown

import "fmt"

func Unmarshal(document *Document, data []byte) error {

	if document == nil {
		return fmt.Errorf("nil document")
	}

	return document.UnmarshalMarkdown(data)

}
