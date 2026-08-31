package yaml

import "fmt"
import net_url "net/url"
import "reflect"

func encodeURL(value reflect.Value) (*Node, error) {

	var raw string

	switch value.Kind() {

	case reflect.Struct:

		parsed := value.Interface().(net_url.URL)
		raw = parsed.String()

	case reflect.Ptr:

		if value.IsNil() {
			return &Node{
				Kind:  ScalarNode,
				Value: "",
			}, nil
		}

		raw = value.Interface().(*net_url.URL).String()

	default:

		return nil, fmt.Errorf(
			"cannot encode %s as net.URL",
			value.Kind(),
		)

	}

	return &Node{
		Kind:  ScalarNode,
		Value: raw,
	}, nil

}
