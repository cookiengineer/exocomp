package yaml

type NodeKind int

const (
	ScalarNode NodeKind = iota
	ObjectNode
	ArrayNode
)

type Node struct {
	Kind           NodeKind
	Value          string
	ObjectChildren map[string]*Node
	ArrayChildren  []*Node
}
