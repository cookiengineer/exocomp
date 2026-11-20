package markdown

import "fmt"

func Marshal(document *Document) ([]byte, error) {

	if document != nil {
		return document.MarshalText()
	} else {
		return []byte{}, fmt.Errorf("nil document")
	}
}
