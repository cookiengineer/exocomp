package types

import "exocomp/schemas"
import utils_api_llamacpp "exocomp/utils/api/llamacpp"
import utils_api_ollama "exocomp/utils/api/ollama"
import utils_api_vllm "exocomp/utils/api/vllm"
import utils_fmt "exocomp/utils/fmt"
import _ "embed"
import "encoding/json"
import "io"
import net_url "net/url"
import "net/http"
import "os"
import "strings"

type Config struct {
	Name        string       `json:"name"`
	Role        string       `json:"role"`
	Model       string       `json:"model"`
	Prompt      string       `json:"prompt"`
	Temperature float64      `json:"temperature"`
	Playground  string       `json:"playground"`
	Sandbox     string       `json:"sandbox"`
	URL         *net_url.URL `json:"url"`
	Debug       bool         `json:"debug"`
}

func NewConfig(name string, role string, model string, prompt string, temperature float64, playground string, sandbox string, url *net_url.URL, debug bool) *Config {

	name   = strings.TrimSpace(name)
	role   = strings.TrimSpace(role)
	model  = strings.TrimSpace(model)
	prompt = utils_fmt.FormatSingleLine(prompt)

	if temperature < 0.0 {
		temperature = 0.0
	} else if temperature > 1.0 {
		temperature = 1.0
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
	}

}

func (config *Config) GetContextLength() int {

	client := &http.Client{}

	request, err1 := http.NewRequest(http.MethodGet, config.ResolveAPI("/v1/models").String(), nil)

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

func (config *Config) ResolveAPI(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}

func (config *Config) MarshalJSON() ([]byte, error) {

	url_str := ""

	if config.URL != nil {
		url_str = config.URL.String()
	}

	return json.Marshal(struct {
		Name        string  `json:"name"`
		Role        string  `json:"role"`
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		Temperature float64 `json:"temperature"`
		Playground  string  `json:"playground"`
		Sandbox     string  `json:"sandbox"`
		URL         string  `json:"url"`
		Debug       bool    `json:"debug"`
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
	})

}

func (config *Config) UnmarshalJSON(data []byte) error {

	var tmp struct {
		Name        string  `json:"name"`
		Role        string  `json:"role"`
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		Temperature float64 `json:"temperature"`
		Playground  string  `json:"playground"`
		Sandbox     string  `json:"sandbox"`
		URL         string  `json:"url"`
		Debug       bool    `json:"debug"`
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
