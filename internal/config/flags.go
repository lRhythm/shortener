package config

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func (t *serverAddress) String() string {
	return string(*t)
}

func (t *serverAddress) Set(v string) error {
	hp := strings.Split(v, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	*t = serverAddress(v)
	return nil
}

func (t *baseURL) String() string {
	return string(*t)
}

func (t *baseURL) Set(v string) error {
	_, err := url.Parse(v)
	if err != nil {
		return err
	}
	*t = baseURL(v)
	return nil
}
