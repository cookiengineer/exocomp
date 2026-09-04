package types

import "encoding/json"
import net_url "net/url"

type Provider struct {
	URL   *net_url.URL `json:"url" yaml:"url"`
	Alias string       `json:"alias" yaml:"alias"`
	Token string       `json:"token" yaml:"token"`
}

func (provider *Provider) MarshalJSON() ([]byte, error) {

	url_str := ""

	if provider.URL != nil {
		url_str = provider.URL.String()
	}

	return json.Marshal(struct {
		URL   string `json:"url"`
		Alias string `json:"alias"`
		Token string `json:"token"`
	}{
		URL:   url_str,
		Alias: provider.Alias,
		Token: provider.Token,
	})

}

func (provider *Provider) UnmarshalJSON(data []byte) error {

	var tmp struct {
		URL   string `json:"url"`
		Alias string `json:"alias"`
		Token string `json:"token"`
	}

	err0 := json.Unmarshal(data, &tmp)

	if err0 == nil {

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
