package containers

import "fmt"

func (tool *Containers) parseArguments(raw any) ([]string, error) {

	result := make([]string, 0)

	if raw == nil {
		return result, nil
	}

	values, ok := raw.([]any)

	if ok == true {

		result = make([]string, len(values))

		for index, value := range values {

			tmp, ok := value.(string)

			if ok == true {
				result[index] = tmp
			} else {
				return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")
			}

		}

		return result, nil

	} else {
		return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")
	}

}
