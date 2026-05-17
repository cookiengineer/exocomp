package vllm

var context_lengths map[string]int

func init() {
	context_lengths = make(map[string]int)
}
