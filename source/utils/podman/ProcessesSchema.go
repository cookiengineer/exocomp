package podman

import "encoding/json"

type ProcessesSchema struct {
	ID     string   `json:"Id"`
	Names  []string `json:"Names"`
	Image  string   `json:"Image"`
	Status string   `json:"Status"`
}

type processNames []string

func (schema *ProcessesSchema) UnmarshalJSON(data []byte) error {

	raw := struct {
		ID     string       `json:"Id"`
		Names  processNames `json:"Names"`
		Image  string       `json:"Image"`
		Status string       `json:"Status"`
	}{}

	err := json.Unmarshal(data, &raw)

	if err == nil {

		schema.ID = raw.ID
		schema.Names = []string(raw.Names)
		schema.Image = raw.Image
		schema.Status = raw.Status

		return nil

	} else {
		return err
	}

}

func (names *processNames) UnmarshalJSON(data []byte) error {

	if len(data) > 0 && data[0] == '"' {

		name := ""

		err := json.Unmarshal(data, &name)

		if err == nil {
			*names = []string{name}
			return nil
		} else {
			return err
		}

	} else {

		values := make([]string, 0)

		err := json.Unmarshal(data, &values)

		if err == nil {
			*names = values
			return nil
		} else {
			return err
		}

	}

}
