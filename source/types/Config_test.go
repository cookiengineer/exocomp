package types

import net_url "net/url"
import "testing"

func TestConfig_ResolvePricing_FromProvider(t *testing.T) {

	config := &Config{
		Providers: map[string]Provider{
			"deepseek-v4-pro:cloud": {
				Model:  "deepseek-v4-pro:cloud",
				Alias:  "deepseek-v4-pro",
				Pricing: Pricing{
					InputPrice:       1.32,
					OutputPrice:      3.96,
					CachedInputPrice: 0.044,
				},
			},
		},
	}

	pricing, ok := config.ResolvePricing("deepseek-v4-pro:cloud")

	if ok != true {
		t.Errorf("Expected pricing to be resolved")
	}

	if pricing.InputPrice != 1.32 {
		t.Errorf("Expected InputPrice %v to be %v", pricing.InputPrice, 1.32)
	}

	if pricing.OutputPrice != 3.96 {
		t.Errorf("Expected OutputPrice %v to be %v", pricing.OutputPrice, 3.96)
	}

	if pricing.CachedInputPrice != 0.044 {
		t.Errorf("Expected CachedInputPrice %v to be %v", pricing.CachedInputPrice, 0.044)
	}

}

func TestConfig_ResolvePricing_Unknown(t *testing.T) {

	config := &Config{
		Providers: map[string]Provider{},
	}

	pricing, ok := config.ResolvePricing("some-model")

	if ok == true {
		t.Errorf("Expected pricing to be unresolved, got %v", pricing)
	}

}

func TestParseProvider_YAML(t *testing.T) {

	data := []byte(`
model: deepseek-v4-pro:cloud
url: https://api.deepseek.com
alias: deepseek-v4-pro
pricing:
  input_price: 1.32
  output_price: 3.96
  cached_input_price: 0.044
`)

	provider, err := ParseProvider(data)

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if provider == nil {
		t.Fatalf("Expected provider to be not nil")
	}

	if provider.Model != "deepseek-v4-pro:cloud" {
		t.Errorf("Expected Model %q to be %q", provider.Model, "deepseek-v4-pro:cloud")
	}

	if provider.Alias != "deepseek-v4-pro" {
		t.Errorf("Expected Alias %q to be %q", provider.Alias, "deepseek-v4-pro")
	}

	if provider.URL == nil {
		t.Errorf("Expected URL to be not nil")
	} else if provider.URL.String() != "https://api.deepseek.com" {
		t.Errorf("Expected URL %q to be %q", provider.URL.String(), "https://api.deepseek.com")
	}

	if provider.Pricing.InputPrice != 1.32 {
		t.Errorf("Expected InputPrice %v to be %v", provider.Pricing.InputPrice, 1.32)
	}

	if provider.Pricing.OutputPrice != 3.96 {
		t.Errorf("Expected OutputPrice %v to be %v", provider.Pricing.OutputPrice, 3.96)
	}

	if provider.Pricing.CachedInputPrice != 0.044 {
		t.Errorf("Expected CachedInputPrice %v to be %v", provider.Pricing.CachedInputPrice, 0.044)
	}

}

func TestParseProvider_JSON(t *testing.T) {

	// NOTE: JSON does not deserialize pricing by design (see Provider.MarshalJSON)
	data := []byte(`{"model":"deepseek-v4-pro:cloud","url":"https://api.deepseek.com","alias":"deepseek-v4-pro","token":"sk-test"}`)

	provider, err := ParseProvider(data)

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if provider == nil {
		t.Fatalf("Expected provider to be not nil")
	}

	if provider.Model != "deepseek-v4-pro:cloud" {
		t.Errorf("Expected Model %q to be %q", provider.Model, "deepseek-v4-pro:cloud")
	}

	if provider.Alias != "deepseek-v4-pro" {
		t.Errorf("Expected Alias %q to be %q", provider.Alias, "deepseek-v4-pro")
	}

	if provider.Token != "sk-test" {
		t.Errorf("Expected Token %q to be %q", provider.Token, "sk-test")
	}

	if provider.URL == nil {
		t.Errorf("Expected URL to be not nil")
	}

	if provider.Pricing != (Pricing{}) {
		t.Errorf("Expected Pricing %v to be empty", provider.Pricing)
	}

}

func TestProvider_MarshalJSON_OmitsPricing(t *testing.T) {

	url, _ := net_url.Parse("https://api.deepseek.com")

	provider := Provider{
		Model: "deepseek-v4-pro:cloud",
		URL:   url,
		Alias: "deepseek-v4-pro",
		Pricing: Pricing{
			InputPrice:  1.32,
			OutputPrice: 3.96,
		},
	}

	data, err := provider.MarshalJSON()

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	parsed, err := ParseProvider(data)

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if parsed == nil {
		t.Fatalf("Expected provider to be not nil")
	}

	if parsed.Pricing != (Pricing{}) {
		t.Errorf("Expected Pricing %v to be empty after JSON round-trip", parsed.Pricing)
	}

	if parsed.Model != "deepseek-v4-pro:cloud" {
		t.Errorf("Expected Model %q to be %q", parsed.Model, "deepseek-v4-pro:cloud")
	}

}
