package yaml

import "reflect"

func encodeSlice(value reflect.Value) (*Node, error) {

	node := &Node{
		Kind:          ArrayNode,
		ArrayChildren: []*Node{},
	}

	for index := 0; index < value.Len(); index++ {

		child, err := encodeValue(value.Index(index))

		if err == nil {
			node.ArrayChildren = append(node.ArrayChildren, child)
		} else {
			return nil, err
		}

	}

	return node, nil

}
