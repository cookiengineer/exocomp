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

		if value.IsNil() == false {

			raw = value.Interface().(*net_url.URL).String()

		} else {
			return &Node{
				Kind:  ScalarNode,
				Value: "",
			}, nil
		}

	default:

		return nil, fmt.Errorf(
			"yaml.EncodeURL: Cannot encode %s as net.URL.",
			value.Kind(),
		)

	}

	return &Node{
		Kind:  ScalarNode,
		Value: raw,
	}, nil

}
