package yaml

type Unmarshaler interface {
	UnmarshalYAML(node *Node) error
}
