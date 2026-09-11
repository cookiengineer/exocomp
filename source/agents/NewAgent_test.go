package agents

import "exocomp/types"
import "testing"

func TestNewAgent_FillsProviderPricing(t *testing.T) {

	config := &types.Config{
		Name:    "Test Agent",
		Role:    "planner",
		Model:   "deepseek-v4-pro:cloud",
		Sandbox: "/tmp/exocomp-test",
		Providers: map[string]types.Provider{
			"deepseek-v4-pro:cloud": {
				Alias: "deepseek-v4-pro",
				Token: "sk-test",
			},
		},
	}

	agent := NewAgent(config)

	if agent == nil {
		t.Fatalf("Expected agent to be not nil")
	}

	provider, ok := config.Providers["deepseek-v4-pro:cloud"]

	if ok != true {
		t.Fatalf("Expected provider to be present")
	}

	if provider.Pricing.InputPrice != 1.32 {
		t.Errorf("Expected InputPrice %v to be %v", provider.Pricing.InputPrice, 1.32)
	}

	if provider.Pricing.OutputPrice != 3.96 {
		t.Errorf("Expected OutputPrice %v to be %v", provider.Pricing.OutputPrice, 3.96)
	}

	if provider.Alias != "deepseek-v4-pro" {
		t.Errorf("Expected Alias %q to be %q", provider.Alias, "deepseek-v4-pro")
	}

	if provider.Token != "sk-test" {
		t.Errorf("Expected Token %q to be %q", provider.Token, "sk-test")
	}

	if provider.URL == nil {
		t.Errorf("Expected URL to be filled from the preset")
	} else if provider.URL.Hostname() != "api.deepseek.com" {
		t.Errorf("Expected URL hostname %q to be %q", provider.URL.Hostname(), "api.deepseek.com")
	}

}

func TestNewAgent_NoProviderPreset(t *testing.T) {

	config := &types.Config{
		Name:    "Test Agent",
		Role:    "planner",
		Model:   "some-local-model",
		Sandbox: "/tmp/exocomp-test",
		Providers: map[string]types.Provider{
			"some-local-model": {
				Alias: "some-local-model",
			},
		},
	}

	agent := NewAgent(config)

	if agent == nil {
		t.Fatalf("Expected agent to be not nil")
	}

	provider, ok := config.Providers["some-local-model"]

	if ok != true {
		t.Fatalf("Expected provider to be present")
	}

	if provider.Pricing != (types.Pricing{}) {
		t.Errorf("Expected Pricing %v to be empty", provider.Pricing)
	}

}
