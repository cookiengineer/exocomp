package yaml

import "reflect"
import "strings"

func Marshal(value any) ([]byte, error) {

    node, err := encodeValue(reflect.ValueOf(value))

    if err != nil {
        return nil, err
    }

    builder := &strings.Builder{}

    writeNode(builder, node, 0)

    return []byte(builder.String()), nil

}

