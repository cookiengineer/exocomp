package markdown

import "fmt"

func Unmarshal(document *Document, data []byte) error {

	if document != nil {
		return document.UnmarshalMarkdown(data)
	} else {
		return fmt.Errorf("markdown.Unmarshal: Nil document.")
	}
}
