package yaml

type NodeKind int

const (
	ScalarNode NodeKind = iota
	ObjectNode
	ArrayNode
)
