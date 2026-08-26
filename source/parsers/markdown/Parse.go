package markdown

import "fmt"

func Parse(file string, bytes []byte) (*Document, error) {

	document := NewDocument(file)

	if document != nil {

		err := document.UnmarshalMarkdown(bytes)

		if err == nil {
			return document, nil
		} else {
			return document, err
		}

	} else {
		return document, fmt.Errorf("Invalid file path \"%s\"", file)
	}

}
