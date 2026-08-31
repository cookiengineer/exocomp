package yaml_test

import "strings"
import "testing"
import "net/url"

import "exocomp/parsers/yaml"

type mockURLProvider struct {
	URL   *url.URL `yaml:"url"`
	Alias string   `yaml:"alias"`
	Token string   `yaml:"token"`
}

type mockURLConfig struct {
	Name      string                     `yaml:"name"`
	URL       *url.URL                   `yaml:"url"`
	Fallback  url.URL                    `yaml:"fallback"`
	Providers map[string]mockURLProvider `yaml:"providers"`
}

func TestUnmarshalURLValueAndPointer(t *testing.T) {

	data := `name: HyRell
url: "https://api.deepseek.com/v1"
fallback: "https://fallback.example.com"
providers:
  deepseek-v4-pro:cloud:
    url: "https://api.deepseek.com"
    alias: "deepseek-v4-pro"
    token: "sk-abc123"
`

	config := mockURLConfig{}

	err := yaml.Unmarshal([]byte(data), &config)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.URL == nil || config.URL.String() != "https://api.deepseek.com/v1" {
		t.Fatalf("expected URL pointer %q, got: %#v", "https://api.deepseek.com/v1", config.URL)
	}

	if config.Fallback.String() != "https://fallback.example.com" {
		t.Fatalf("expected Fallback value %q, got: %#v", "https://fallback.example.com", config.Fallback)
	}

	provider, exists := config.Providers["deepseek-v4-pro:cloud"]

	if exists == false {
		t.Fatalf("expected provider to exist, got: %#v", config.Providers)
	}

	if provider.URL == nil || provider.URL.String() != "https://api.deepseek.com" {
		t.Fatalf("expected provider URL %q, got: %#v", "https://api.deepseek.com", provider.URL)
	}

}

func TestUnmarshalURLInvalidString(t *testing.T) {

	data := "url: \"://invalid\"\n"

	config := struct {
		URL *url.URL `yaml:"url"`
	}{}

	err := yaml.Unmarshal([]byte(data), &config)

	if err == nil {
		t.Fatalf("expected an error for invalid url, got nil")
	}

}

func TestMarshalURLPointer(t *testing.T) {

	parsed, err := url.Parse("https://example.com/v1")

	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}

	value := struct {
		URL *url.URL `yaml:"url"`
	}{
		URL: parsed,
	}

	out, err := yaml.Marshal(value)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if strings.Contains(string(out), "url: https://example.com/v1") == false {
		t.Fatalf("expected marshaled url, got: %q", string(out))
	}

}

func TestMarshalURLValue(t *testing.T) {

	parsed, err := url.Parse("https://example.com/v1")

	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}

	value := struct {
		URL url.URL `yaml:"url"`
	}{
		URL: *parsed,
	}

	out, err := yaml.Marshal(value)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if strings.Contains(string(out), "url: https://example.com/v1") == false {
		t.Fatalf("expected marshaled url, got: %q", string(out))
	}

}

func TestMarshalURLNilPointer(t *testing.T) {

	value := struct {
		URL *url.URL `yaml:"url"`
	}{}

	out, err := yaml.Marshal(value)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if strings.Contains(string(out), "url:") == false {
		t.Fatalf("expected marshaled url key, got: %q", string(out))
	}

}
