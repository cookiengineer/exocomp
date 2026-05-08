package types

type Tool interface {
	Call(string, map[string]interface{}) (string, error)
	Get(string)                          (any, error)
}
