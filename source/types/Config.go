package types

import "exocomp/schemas"
import utils_api_llamacpp "exocomp/utils/api/llamacpp"
import utils_api_ollama "exocomp/utils/api/ollama"
import utils_api_vllm "exocomp/utils/api/vllm"
import utils_fmt "exocomp/utils/fmt"
import "exocomp/parsers/yaml"
import _ "embed"
import "encoding/json"
import "fmt"
import "io"
import net_url "net/url"
import "net/http"
import "os"
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

func NewConfig(name string, role string, model string, prompt string, temperature float64, playground string, sandbox string, url *net_url.URL, debug bool) *Config {

	if role == "planner" && GlobalConfig != nil {

		if name == "" {
			name = GlobalConfig.Name
		}

		if model == "" {
			model = GlobalConfig.Model
		}

		if prompt == "" {
			prompt = GlobalConfig.Prompt
		}

		if temperature == 0.0 {
			temperature = GlobalConfig.Temperature
		}

		if url == nil {
			url = GlobalConfig.URL
		}

		if debug == false {
			debug = GlobalConfig.Debug
		}

	}

	name   = strings.TrimSpace(name)
	role   = strings.TrimSpace(role)
	model  = strings.TrimSpace(model)
	prompt = utils_fmt.FormatSingleLine(prompt)

	if temperature < 0.0 {
		temperature = 0.0
	} else if temperature > 1.0 {
		temperature = 1.0
	}

	if role == "" {
		role = "planner"
	}

	if playground == "" {

		base := os.TempDir()
		tmp, err := os.MkdirTemp(base, "exocomp-playground-*")

		if err == nil {
			playground = tmp
		} else {
			playground = "/tmp/exocomp"
		}

	}

	if sandbox == "" {

		cwd, err := os.Getwd()

		if err == nil {
			sandbox = cwd
		} else {
			sandbox = "/tmp/exocomp/sandbox"
		}

	}

	if url == nil {

		tmp, err := net_url.Parse("http://localhost:11434/v1")

		if err == nil {
			url = tmp
		}

	}

	providers := make(map[string]Provider)

	if GlobalConfig != nil {

		if role == "planner" {

			for model, provider := range GlobalConfig.Providers {
				providers[model] = provider
			}

		} else if model != "" {

			provider, ok := GlobalConfig.Providers[model]

			if ok == true {
				providers[model] = provider
			}

		}

	}

	return &Config{
		Name:        name,
		Role:        role,
		Model:       model,
		Prompt:      prompt,
		Temperature: temperature,
		Playground:  playground,
		Sandbox:     sandbox,
		URL:         url,
		Debug:       debug,
		Providers:   providers,
	}

}

func ParseConfig(data []byte) (*Config, error) {

	if len(data) > 2 && data[0] == '{' && data[len(data)-1] == '}' {

		config := Config{}
		err    := json.Unmarshal(data, &config)

		if err == nil {
			return &config, nil
		} else {
			return nil, err
		}

	} else {

		config := Config{}
		err    := yaml.Unmarshal(data, &config)

		if err == nil {
			return &config, nil
		} else {
			return nil, err
		}

	}

}

func (config *Config) GetContextLength(model string) int {

	client       := &http.Client{}
	resolved_url := config.ResolveURL(model, "/models")

	request, err1 := http.NewRequest(http.MethodGet, resolved_url.String(), nil)

	if err1 == nil {

		response, err2 := client.Do(request)

		if err2 == nil {

			response_payload, err3 := io.ReadAll(response.Body)

			if err3 == nil {

				schema := schemas.ModelsResponse{}
				err4   := json.Unmarshal(response_payload, &schema)

				if err4 == nil {

					server_type := schema.OwnedBy()

					if server_type == "llamacpp" {

						return utils_api_llamacpp.GetContextLength(config.URL, config.Model)

					} else if server_type == "ollama" {

						return utils_api_ollama.GetContextLength(config.URL, config.Model)

					} else if server_type == "vllm" {

						return utils_api_vllm.GetContextLength(config.URL, config.Model)

					}

				}

			}

		}

	}

	return 0

}

func (config *Config) GetPrompt() string {
	return strings.TrimSpace(config.Prompt)
}

func (config *Config) ResolveToken(model string) string {

	provider, ok := config.Providers[model]

	if ok == true {
		return strings.TrimSpace(provider.Token)
	} else {
		return ""
	}

}

func (config *Config) ResolveModel(model string) string {

	provider, ok := config.Providers[model]

	if ok == true {

		if provider.Alias != "" {
			return provider.Alias
		} else {
			return model
		}

	} else {
		return model
	}

}

func (config *Config) ResolveURL(model string, path string) *net_url.URL {

	base_url := config.URL
	api_path := ""

	provider, ok := config.Providers[model]

	if ok == true {
		base_url, _ = net_url.Parse(provider.URL.String())
	}

	if strings.HasPrefix(base_url.Path, "/") && len(base_url.Path) > 1 {

		// "/v1" or "/v1/"
		tmp_base := base_url.Path

		if strings.HasSuffix(tmp_base, "/") {
			tmp_base = strings.TrimSpace(tmp_base[0:len(tmp_base)-1])
		}

		// "/chat/completions"
		tmp_path := path

		if strings.HasPrefix(tmp_path, "/") {
			tmp_path = strings.TrimSpace(tmp_path[1:])
		}

		// "/v1/chat/completions"
		api_path = fmt.Sprintf("%s/%s", tmp_base, tmp_path)

	} else if strings.HasPrefix(path, "/") {
		api_path = strings.TrimSpace(path)
	}

	if api_path != "" {

		return base_url.ResolveReference(&net_url.URL{
			Path: api_path,
		})

	} else {
		return base_url
	}

}

func (config *Config) Public() *Config {

	clone := *config
	clone.Providers = make(map[string]Provider)

	return &clone

}

func (config *Config) MarshalJSON() ([]byte, error) {

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

	if err0 == nil {

		config.Name        = tmp.Name
		config.Role        = tmp.Role
		config.Model       = tmp.Model
		config.Prompt      = tmp.Prompt
		config.Temperature = tmp.Temperature
		config.Playground  = tmp.Playground
		config.Sandbox     = tmp.Sandbox
		config.Debug       = tmp.Debug
		config.Providers   = tmp.Providers

		tmp_url, err1 := net_url.Parse(tmp.URL)

		if err1 == nil {
			config.URL = tmp_url
		}

		return nil

	} else {
		return err0
	}

}

func (config *Config) MarshalYAML() ([]byte, error) {

	url_str := ""

	if config.URL != nil {
		url_str = config.URL.String()
	}

	return yaml.Marshal(struct {
		Name        string              `yaml:"name"`
		Role        string              `yaml:"role"`
		Model       string              `yaml:"model"`
		Prompt      string              `yaml:"prompt"`
		Temperature float64             `yaml:"temperature"`
		Playground  string              `yaml:"playground"`
		Sandbox     string              `yaml:"sandbox"`
		URL         string              `yaml:"url"`
		Debug       bool                `yaml:"debug"`
		Providers   map[string]Provider `yaml:"providers"`
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

func (config *Config) UnmarshalYAML(data []byte) error {

	var tmp struct {
		Name        string              `yaml:"name"`
		Role        string              `yaml:"role"`
		Model       string              `yaml:"model"`
		Prompt      string              `yaml:"prompt"`
		Temperature float64             `yaml:"temperature"`
		Playground  string              `yaml:"playground"`
		Sandbox     string              `yaml:"sandbox"`
		URL         string              `yaml:"url"`
		Debug       bool                `yaml:"debug"`
		Providers   map[string]Provider `yaml:"providers"`
	}

	err0 := yaml.Unmarshal(data, &tmp)

	if err0 == nil {

		config.Name        = tmp.Name
		config.Role        = tmp.Role
		config.Model       = tmp.Model
		config.Prompt      = tmp.Prompt
		config.Temperature = tmp.Temperature
		config.Playground  = tmp.Playground
		config.Sandbox     = tmp.Sandbox
		config.Debug       = tmp.Debug
		config.Providers   = tmp.Providers

		tmp_url, err1 := net_url.Parse(tmp.URL)

		if err1 == nil {
			config.URL = tmp_url
		}

		return nil

	} else {
		return err0
	}

}

func (config *Config) Update(name string, role string, model string, prompt string, temperature float64) {

	prompt = utils_fmt.FormatSingleLine(prompt)

	if temperature < 0.1 {
		temperature = 0.1
	} else if temperature > 1.0 {
		temperature = 1.0
	}

	config.Name        = name
	config.Role        = role
	config.Model       = model
	config.Prompt      = prompt
	config.Temperature = temperature

}
