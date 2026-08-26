package html

import "fmt"

func Parse(file string, bytes []byte) (*Document, error) {

	document := NewDocument(file)

	if document != nil {

		err := document.Parse(bytes)

		if err == nil {
			return document, nil
		} else {
			return document, err
		}

	} else {
		return nil, fmt.Errorf("Invalid URL \"%s\"", file)
	}

}
