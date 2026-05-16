package llamacpp

type ModelMeta struct {
	ContextTrainingLength int `json:"n_ctx_train"`
	EmbeddingLength       int `json:"n_embd"`
	Parameters            int `json:"n_params"`
	Size                  int `json:"size"`
	Vocabulary            int `json:"n_vocab"`
	VocabularyType        int `json:"vocab_type"`
}
