package types

// Pricing describes the cost of a model in US Dollars per 1M (one million)
// tokens. CachedInputPrice is the discounted price for input tokens that hit
// the provider's prompt cache; when it is zero the cache-miss InputPrice is
// used for all input tokens.
type Pricing struct {
	InputPrice       float64 `json:"input_price" yaml:"input_price"`
	OutputPrice      float64 `json:"output_price" yaml:"output_price"`
	CachedInputPrice float64 `json:"cached_input_price" yaml:"cached_input_price"`
}

type Provider struct {
	URL     string  `json:"url" yaml:"url"`
	Alias   string  `json:"alias" yaml:"alias"`
	Token   string  `json:"token" yaml:"token"`
	Pricing Pricing `json:"pricing" yaml:"pricing"`
}
