package config

import (
	"net/url"
)

type serverAddress string
type baseURL string
type fileStoragePath string
type databaseDSN string

type Cfg struct {
	ServerAddress   serverAddress   `env:"SERVER_ADDRESS"`
	BaseURL         baseURL         `env:"BASE_URL"`
	FileStoragePath fileStoragePath `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     databaseDSN     `env:"DATABASE_DSN"`
	CookieSecretKey string          `env:"COOKIE_SECRET_KEY" envDefault:"o04n+9H6PWZs8PSxQqh9R1bWDL3sEUMfzx1gg0XTWns="`
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

func (c *Cfg) DSN() (string, bool) {
	dsn := string(c.DatabaseDSN)
	return dsn, len(dsn) > 0
}

func (c *Cfg) CookieKey() string {
	return c.CookieSecretKey
}
