package types

import "encoding/json"
import net_url "net/url"
import "strings"

type Config struct {
	Name        string              `json:"name" yaml:"name"`
	Role        string              `json:"role" yaml:"role"`
	Model       string              `json:"model" yaml:"model"`
	Prompt      string              `json:"prompt" yaml:"prompt"`
	Temperature float64             `json:"temperature" yaml:"temperature"`
	Playground  string              `json:"playground" yaml:"playground"`
	Sandbox     string              `json:"sandbox" yaml:"sandbox"`
	URL         *net_url.URL        `json:"url" yaml:"url"`
	Debug       bool                `json:"debug" yaml:"debug"`
	Providers   map[string]Provider `json:"providers" yaml:"providers"`
}

func NewConfig() *Config {

	url, _ := net_url.Parse("http://localhost:3000/")

	config := Config{
		URL:       url,
		Providers: make(map[string]Provider),
	}

	return &config

}

func (config Config) MarshalJSON() ([]byte, error) {

	url_str := ""

	if config.URL != nil {
		url_str = config.URL.String()
	}

	return json.Marshal(struct {
		Name        string              `json:"name"`
		Role        string              `json:"role"`
		Model       string              `json:"model"`
		Prompt      string              `json:"prompt"`
		Temperature float64             `json:"temperature"`
		Playground  string              `json:"playground"`
		Sandbox     string              `json:"sandbox"`
		URL         string              `json:"url"`
		Debug       bool                `json:"debug"`
		Providers   map[string]Provider `json:"providers,omitempty"`
	}{
		Name:        config.Name,
		Role:        config.Role,
		Model:       config.Model,
		Prompt:      config.Prompt,
		Temperature: config.Temperature,
		Playground:  config.Playground,
		Sandbox:     config.Sandbox,
		URL:         url_str,
		Debug:       config.Debug,
		Providers:   config.Providers,
	})

}

func (config *Config) UnmarshalJSON(data []byte) error {

	var tmp struct {
		Name        string              `json:"name"`
		Role        string              `json:"role"`
		Model       string              `json:"model"`
		Prompt      string              `json:"prompt"`
		Temperature float64             `json:"temperature"`
		Playground  string              `json:"playground"`
		Sandbox     string              `json:"sandbox"`
		URL         string              `json:"url"`
		Debug       bool                `json:"debug"`
		Providers   map[string]Provider `json:"providers"`
	}

	err0 := json.Unmarshal(data, &tmp)

	if err0 != nil {
		return err0
	}

	config.Name = tmp.Name
	config.Role = tmp.Role
	config.Model = tmp.Model
	config.Prompt = tmp.Prompt
	config.Temperature = tmp.Temperature
	config.Playground = tmp.Playground
	config.Sandbox = tmp.Sandbox
	config.Debug = tmp.Debug
	config.Providers = tmp.Providers

	tmp_url, err1 := net_url.Parse(tmp.URL)

	if err1 == nil {
		config.URL = tmp_url
	}

	return nil

}

func (config *Config) GetPrompt() string {
	return strings.TrimSpace(config.Prompt)
}

func (config *Config) ResolveAPI(path string) string {

	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path

}
