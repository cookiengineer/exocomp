package types

import "exocomp/parsers/yaml"
import "encoding/json"
import net_url "net/url"

type Provider struct {
	Model   string       `json:"model" yaml:"model"`
	URL     *net_url.URL `json:"url" yaml:"url"`
	Alias   string       `json:"alias" yaml:"alias"`
	Token   string       `json:"token" yaml:"token"`
	Pricing Pricing      `json:"pricing" yaml:"pricing"`
}

func ParseProvider(data []byte) (*Provider, error) {

	if len(data) > 2 && data[0] == '{' && data[len(data)-1] == '}' {

		provider := Provider{}
		err := json.Unmarshal(data, &provider)

		if err == nil {
			return &provider, nil
		} else {
			return nil, err
		}

	} else {

		provider := Provider{}
		err := yaml.Unmarshal(data, &provider)

		if err == nil {
			return &provider, nil
		} else {
			return nil, err
		}

	}

}

func (provider Provider) MarshalJSON() ([]byte, error) {

	url_str := ""

	if provider.URL != nil {
		url_str = provider.URL.String()
	}

	return json.Marshal(struct {
		Model string `json:"model,omitempty"`
		URL   string `json:"url"`
		Alias string `json:"alias"`
		Token string `json:"token"`
	}{
		Model: provider.Model,
		URL:   url_str,
		Alias: provider.Alias,
		Token: provider.Token,
	})

}

func (provider *Provider) UnmarshalJSON(data []byte) error {

	var tmp struct {
		Model string `json:"model"`
		URL   string `json:"url"`
		Alias string `json:"alias"`
		Token string `json:"token"`
	}

	err0 := json.Unmarshal(data, &tmp)

	if err0 == nil {

		provider.Model = tmp.Model
		provider.Alias = tmp.Alias
		provider.Token = tmp.Token

		tmp_url, err1 := net_url.Parse(tmp.URL)

		if err1 == nil {
			provider.URL = tmp_url
		}

		return nil

	} else {
		return err0
	}

}
