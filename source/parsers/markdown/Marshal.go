package markdown

import "fmt"

func Marshal(document *Document) ([]byte, error) {

	if document == nil {
		return []byte{}, fmt.Errorf("nil document")
	}

	return document.MarshalText()

}
