package config

import (
	"errors"
	"flag"
	"net/url"
	"strconv"
	"strings"
)

type Cfg struct {
	address
	path
}

type address struct {
	host string
	port int
}

type path struct {
	prefix string
}

func (c *Cfg) Host() string {
	return c.address.String()
}

func (c *Cfg) Path() string {
	return c.path.String()
}

func (a *address) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

func (a *address) Set(v string) error {
	hp := strings.Split(v, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.host = hp[0]
	a.port = port
	return nil
}

func (p *path) String() string {
	return p.prefix
}

func (p *path) Set(v string) error {
	u, err := url.Parse(v)
	if err != nil {
		return err
	}
	p.prefix = u.Path
	return nil
}

func New() *Cfg {
	a := &address{
		host: "localhost",
		port: 8080,
	}
	p := &path{
		prefix: "",
	}
	if flag.Lookup("a") == nil {
		_ = flag.Value(a)
		flag.Var(a, "a", "Net address host:port")
	}
	if flag.Lookup("b") == nil {
		_ = flag.Value(p)
		flag.Var(p, "b", "Net address with route prefix (example: http://localhost:8080/prefix)")
	}
	flag.Parse()
	return &Cfg{
		address: *a,
		path:    *p,
	}
}
