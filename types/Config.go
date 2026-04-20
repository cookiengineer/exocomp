package types

import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import "bytes"
import _ "embed"
import "encoding/json"
import "io"
import net_url "net/url"
import "net/http"
import "os"
import "strings"

type Config struct {
	Name        string
	Agent       string
	Debug       bool
	Model       string
	Playground  string
	Prompt      string
	Sandbox     string
	Temperature float64
	URL         *net_url.URL
}

func NewConfig(name string, agent string, debug bool, model string, playground string, prompt string, sandbox string, temperature float64, url *net_url.URL) *Config {

	prompt = utils_fmt.FormatSingleLine(prompt)

	if temperature < 0.1 {
		temperature = 0.1
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

	return &Config{
		Name:        name,
		Agent:       agent,
		Debug:       debug,
		Model:       model,
		Playground:  playground,
		Prompt:      prompt,
		Sandbox:     sandbox,
		Temperature: temperature,
		URL:         url,
	}

}

func (config *Config) GetContextLength() int {

	// TODO: This is ollama specific and won't
	// work with only OpenAI compatible APIs

	client := &http.Client{}

	request_body, _ := json.Marshal(schemas.ShowRequest{
		Name: config.Model,
	})

	request, err1 := http.NewRequest(http.MethodPost, config.ResolveAPI("/api/show").String(), bytes.NewBuffer(request_body))

	if err1 == nil {

		request.Header.Set("Content-Type", "application/json")

		response, err2 := client.Do(request)

		if err2 == nil {

			response_payload, err3 := io.ReadAll(response.Body)

			if err3 == nil {

				schema := schemas.ShowResponse{}
				err4 := json.Unmarshal(response_payload, &schema)

				if err4 == nil {
					return schema.ContextLength()
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
