package types

import net_url "net/url"

type Provider struct {
	URL   *net_url.URL `json:"url" yaml:"url"`
	Alias string       `json:"alias" yaml:"alias"`
	Token string       `json:"token" yaml:"token"`
}

