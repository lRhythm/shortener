package config

import "net/url"

type (
	// serverAddress - тип для адреса запуска сервиса.
	serverAddress string
	// baseURL - тип для base URL сервиса.
	baseURL string
	// fileStoragePath - тип для пути файла для хранения данных.
	fileStoragePath string
	// databaseDSN - тип для PostgreSQL DSN.
	databaseDSN string
)

// Cfg - структура описывающая конфигурацию сервиса.
type Cfg struct {
	ServerAddress   serverAddress   `env:"SERVER_ADDRESS"`
	BaseURL         baseURL         `env:"BASE_URL"`
	FileStoragePath fileStoragePath `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     databaseDSN     `env:"DATABASE_DSN"`
	CookieSecretKey string          `env:"COOKIE_SECRET_KEY" envDefault:"o04n+9H6PWZs8PSxQqh9R1bWDL3sEUMfzx1gg0XTWns="`
}

// Host - получение адреса запуска сервиса.
func (c *Cfg) Host() string {
	return string(c.ServerAddress)
}

// Path - получение значения base URL сервиса.
func (c *Cfg) Path() string {
	u, err := url.ParseRequestURI(string(c.BaseURL))
	if err != nil {
		return ""
	}
	return u.Path
}

// File - получение пути файла для хранения данных.
func (c *Cfg) File() string {
	return string(c.FileStoragePath)
}

// DSN - получение PostgreSQL DSN.
func (c *Cfg) DSN() (string, bool) {
	dsn := string(c.DatabaseDSN)
	return dsn, len(dsn) > 0
}

// CookieKey - получение ключа шифрования cookie.
func (c *Cfg) CookieKey() string {
	return c.CookieSecretKey
}
