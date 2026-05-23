package yaml

import "reflect"

func decodeSlice(node *Node, target reflect.Value) error {

	target_type := target.Type().Elem()
	slice       := reflect.MakeSlice(
		target.Type(),
		0,
		len(node.ArrayChildren),
	)

	for _, child := range node.ArrayChildren {

		element := reflect.New(target_type).Elem()
		err     := decodeValue(child, element)

		if err == nil {
			slice = reflect.Append(slice, element)
		} else {
			return err
		}

	}

	target.Set(slice)

	return nil

}
