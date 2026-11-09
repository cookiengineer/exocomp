package yaml

type Node struct {
	Kind           NodeKind
	Value          string
	ObjectChildren map[string]*Node
	ArrayChildren  []*Node
}
