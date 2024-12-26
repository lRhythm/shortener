package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

// New - создание и заполнение структуры конфигурации сервиса.
func New() (*Cfg, error) {
	var cfg Cfg
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	c, err := cfg.withFlags().withConfig()
	if err != nil {
		return nil, err
	}
	return c.withDefault(), nil
}

// withFlags - поддержка флагов для заполнения config.
func (c *Cfg) withFlags() *Cfg {
	sa := new(serverAddress)
	bu := new(baseURL)
	fsp := new(fileStoragePath)
	dd := new(databaseDSN)
	var tls bool
	var config string
	var needParse bool
	if flag.Lookup("a") == nil {
		_ = flag.Value(sa)
		flag.Var(sa, "a", "Net address host:port")
		needParse = true
	}
	if flag.Lookup("b") == nil {
		_ = flag.Value(bu)
		flag.Var(bu, "b", "Net address with route prefix (example: http://localhost:8080/prefix)")
		needParse = true
	}
	if flag.Lookup("f") == nil {
		_ = flag.Value(fsp)
		flag.Var(fsp, "f", "File storage path (example: ./storage")
		needParse = true
	}
	if flag.Lookup("d") == nil {
		_ = flag.Value(dd)
		flag.Var(dd, "d", "PostgreSQL DSN")
		needParse = true
	}
	if flag.Lookup("s") == nil {
		flag.BoolVar(&tls, "s", false, "Enable HTTPS (TLS)")
		needParse = true
	}
	if flag.Lookup("c") == nil {
		flag.StringVar(&config, "c", "", "Config file path (example: ./configs/config.json)")
		needParse = true
	}
	if flag.Lookup("config") == nil {
		flag.StringVar(&config, "config", "", "Config file path (example: ./configs/config.json)")
		needParse = true
	}
	if needParse {
		flag.Parse()
	}
	if *sa != "" && c.ServerAddress == "" {
		c.ServerAddress = *sa
	}
	if *bu != "" && c.BaseURL == "" {
		c.BaseURL = *bu
	}
	if *fsp != "" && c.FileStoragePath == "" {
		c.FileStoragePath = *fsp
	}
	if *dd != "" && c.DatabaseDSN == "" {
		c.DatabaseDSN = *dd
	}
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "c" || f.Name == "config" {
			c.Config = config
		}
	})
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "s" {
			*c.EnableHTTPS = tls
		}
	})
	return c
}

// withConfig - установка значений полей в Cfg из config файла (*.json), если путь к файлу передан как флаг или env.
func (c *Cfg) withConfig() (*Cfg, error) {
	if path := c.Config; len(c.Config) > 0 {
		var config Cfg
		configFile, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			return nil, err
		}
		if len(c.ServerAddress) == 0 && len(config.ServerAddress) > 0 {
			c.ServerAddress = config.ServerAddress
		}
		if len(c.BaseURL) == 0 && len(config.BaseURL) > 0 {
			c.BaseURL = config.BaseURL
		}
		if len(c.FileStoragePath) == 0 && len(config.FileStoragePath) > 0 {
			c.FileStoragePath = config.FileStoragePath
		}
		if len(c.DatabaseDSN) == 0 && len(config.DatabaseDSN) > 0 {
			c.DatabaseDSN = config.DatabaseDSN
		}
		if c.EnableHTTPS == nil && *config.EnableHTTPS {
			c.EnableHTTPS = config.EnableHTTPS
		}
	}
	return c, nil
}

// withDefault - установка значений полей в config по умолчанию, если значения не определены ранее.
func (c *Cfg) withDefault() *Cfg {
	if c.ServerAddress == "" {
		c.ServerAddress = "localhost:8080"
	}
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost"
	}
	if c.FileStoragePath == "" {
		c.FileStoragePath = "./storage"
	}
	if c.EnableHTTPS == nil {
		var b bool
		c.EnableHTTPS = &b
	}
	return c
}
