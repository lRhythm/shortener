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
	if c.ServerAddress == "" {
		_ = flag.Value(sa)
		flag.Var(sa, "a", "Net address host:port")
	}
	if c.BaseURL == "" {
		_ = flag.Value(bu)
		flag.Var(bu, "b", "Net address with route prefix (example: http://localhost:8080/prefix)")
	}
	if c.ServerAddress == "" || c.BaseURL == "" {
		flag.Parse()
	}
	if *sa != "" {
		c.ServerAddress = *sa
	}
	if *bu != "" {
		c.BaseURL = *bu
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
	return c
}
