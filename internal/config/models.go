package config

import (
	"net/url"
)

type serverAddress string
type baseURL string

type Cfg struct {
	ServerAddress serverAddress `env:"SERVER_ADDRESS"`
	BaseURL       baseURL       `env:"BASE_URL"`
}

func (c *Cfg) Host() string {
	return string(c.ServerAddress)
}

func (c *Cfg) Path() string {
	u, err := url.Parse(string(c.BaseURL))
	if err != nil {
		return ""
	}
	return u.Path
}
