package vllm

var contextLengths map[string]int

func init() {
	contextLengths = make(map[string]int)
}
