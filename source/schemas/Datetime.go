package schemas

import "exocomp/encoding/yaml"
import "encoding/json"
import "fmt"
import "time"

type Datetime struct {
	time.Time
}

func NewDatetime() Datetime {

	return Datetime{
		Time: time.Now(),
	}

}

func (datetime Datetime) MarshalJSON() ([]byte, error) {
	return json.Marshal(datetime.Format("2006-01-02 15:04:05"))
}

func (datetime *Datetime) UnmarshalJSON(data []byte) error {

	var datetime_str string

	err0 := json.Unmarshal(data, &datetime_str)

	if err0 == nil {

		tmp1, err1 := time.Parse("2006-01-02 15:04:05", datetime_str)

		if err1 == nil {

			datetime.Time = tmp1

			return nil

		} else {
			return fmt.Errorf("invalid datetime format: %w", err1)
		}

	} else {
		return err0
	}

}

func (datetime Datetime) MarshalYAML() (*yaml.Node, error) {

	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: datetime.Format("2006-01-02 15:04:05"),
	}, nil

}

func (datetime *Datetime) UnmarshalYAML(node *yaml.Node) error {

	if node.Kind == yaml.ScalarNode {

		tmp1, err1 := time.Parse("2006-01-02 15:04:05", node.Value)

		if err1 == nil {

			datetime.Time = tmp1

			return nil

		} else {
			return fmt.Errorf("invalid datetime format: %w", err1)
		}

	} else {
		return fmt.Errorf("invalid yaml datetime format: expected ScalarNode, got %v", node.Kind)
	}

}
