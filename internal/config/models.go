package config

import (
	"net/url"
)

type serverAddress string
type baseURL string
type fileStoragePath string

type Cfg struct {
	ServerAddress   serverAddress   `env:"SERVER_ADDRESS"`
	BaseURL         baseURL         `env:"BASE_URL"`
	FileStoragePath fileStoragePath `env:"FILE_STORAGE_PATH"`
}

func (c *Cfg) Host() string {
	return string(c.ServerAddress)
}

func (c *Cfg) Path() string {
	u, err := url.ParseRequestURI(string(c.BaseURL))
	if err != nil {
		return ""
	}
	return u.Path
}

func (c *Cfg) File() string {
	return string(c.FileStoragePath)
}
