package yaml

type Marshaler interface {
	MarshalYAML() (*Node, error)
}
