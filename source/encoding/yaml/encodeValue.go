package yaml

import "fmt"
import "reflect"

func encodeValue(value reflect.Value) (*Node, error) {

	if marshaler, ok := value.Interface().(Marshaler); ok == true {
		return marshaler.MarshalYAML()
	}

	switch value.Kind() {

	case reflect.Bool:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%t", value.Bool()),
		}, nil

	case reflect.Float32:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%f", value.Float()),
		}, nil

	case reflect.Float64:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%f", value.Float()),
		}, nil

	case reflect.Int:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Int()),
		}, nil

	case reflect.Int8:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Int()),
		}, nil

	case reflect.Int16:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Int()),
		}, nil

	case reflect.Int32:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Int()),
		}, nil

	case reflect.Int64:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Int()),
		}, nil

	case reflect.Uint:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Uint()),
		}, nil

	case reflect.Uint8:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Uint()),
		}, nil

	case reflect.Uint16:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Uint()),
		}, nil

	case reflect.Uint32:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Uint()),
		}, nil

	case reflect.Uint64:

		return &Node{
			Kind:  ScalarNode,
			Value: fmt.Sprintf("%d", value.Uint()),
		}, nil

	case reflect.Slice:

		return encodeSlice(value)

	case reflect.String:

		return &Node{
			Kind:  ScalarNode,
			Value: value.String(),
		}, nil

	case reflect.Struct:

		return encodeStruct(value)

	default:

		return nil, fmt.Errorf(
			"unsupported type: %s",
			value.Kind(),
		)

	}
}
