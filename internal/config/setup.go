package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

func New() (*Cfg, error) {
	var cfg Cfg
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg.withFlags().withDefault(), nil
}

func (c *Cfg) withFlags() *Cfg {
	sa := new(serverAddress)
	bu := new(baseURL)
	fsp := new(fileStoragePath)
	dd := new(databaseDSN)
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
	return c
}

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
	return c
}
