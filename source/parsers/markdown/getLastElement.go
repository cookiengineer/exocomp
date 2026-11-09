package markdown

func getLastElement(result []*Element, typ string) *Element {

	var pointer *Element = nil

	if len(result) > 0 {

		tmp := result[len(result)-1]

		if tmp.Type == typ {
			pointer = tmp
		}

	}

	return pointer

}
