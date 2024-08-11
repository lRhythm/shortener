package config

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func (t *serverAddress) UnmarshalText(text []byte) error {
	v := string(text)
	if v == "" {
		return nil
	}
	hp := strings.Split(string(text), ":")
	if len(hp) != 2 {
		return errors.New("the environment variable \"SERVER_ADDRESS\" value must be format host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return errors.New("the environment variable \"SERVER_ADDRESS\" value :port must be an integer")
	}
	*t = serverAddress(v)
	return nil
}

func (t *baseURL) UnmarshalText(text []byte) error {
	v := string(text)
	if v == "" {
		return nil
	}
	_, e := url.ParseRequestURI(v)
	if e != nil {
		return errors.New("the environment variable \"BASE_URL\" value must be a valid URL")
	}
	*t = baseURL(v)
	return nil
}
