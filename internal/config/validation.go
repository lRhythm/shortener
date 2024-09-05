package config

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (t *serverAddress) validate(v string) error {
	hp := strings.Split(v, ":")
	if len(hp) != 2 {
		return errors.New("need address in a format host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	return nil
}

func (t *baseURL) validate(v string) error {
	_, err := url.ParseRequestURI(v)
	return err
}

func (t *fileStoragePath) validate(v string) error {
	f, err := os.OpenFile(v, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
