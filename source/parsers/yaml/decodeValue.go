package yaml

import "fmt"
import "reflect"
import "strconv"

func decodeValue(node *Node, target reflect.Value) error {

	if target.CanAddr() {

		unmarshaler, implements := target.Addr().Interface().(Unmarshaler)

		if implements == true {
			return unmarshaler.UnmarshalYAML(node)
		}

	}

	switch target.Kind() {

	case reflect.Bool:

		tmp, err := strconv.ParseBool(node.Value)

		if err == nil {
			target.SetBool(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Int:

		tmp, err := strconv.ParseInt(node.Value, 10, 64)

		if err == nil {
			target.SetInt(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Int8:

		tmp, err := strconv.ParseInt(node.Value, 10, 8)

		if err == nil {
			target.SetInt(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Int16:

		tmp, err := strconv.ParseInt(node.Value, 10, 16)

		if err == nil {
			target.SetInt(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Int32:

		tmp, err := strconv.ParseInt(node.Value, 10, 32)

		if err == nil {
			target.SetInt(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Int64:

		tmp, err := strconv.ParseInt(node.Value, 10, 64)

		if err == nil {
			target.SetInt(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Uint:

		tmp, err := strconv.ParseUint(node.Value, 10, 64)

		if err == nil {
			target.SetUint(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Uint8:

		tmp, err := strconv.ParseUint(node.Value, 10, 8)

		if err == nil {
			target.SetUint(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Uint16:

		tmp, err := strconv.ParseUint(node.Value, 10, 16)

		if err == nil {
			target.SetUint(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Uint32:

		tmp, err := strconv.ParseUint(node.Value, 10, 32)

		if err == nil {
			target.SetUint(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Uint64:

		tmp, err := strconv.ParseUint(node.Value, 10, 64)

		if err == nil {
			target.SetUint(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Float32:

		tmp, err := strconv.ParseFloat(node.Value, 32)

		if err == nil {
			target.SetFloat(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Float64:

		tmp, err := strconv.ParseFloat(node.Value, 64)

		if err == nil {
			target.SetFloat(tmp)
			return nil
		} else {
			return err
		}

	case reflect.Map:

		return decodeMap(node, target)

	case reflect.Slice:

		return decodeSlice(node, target)

	case reflect.String:

		target.SetString(node.Value)
		return nil

	case reflect.Struct:

		return decodeStruct(node, target)

	default:

		return fmt.Errorf(
			"unsupported type: %s",
			target.Kind(),
		)

	}

}
