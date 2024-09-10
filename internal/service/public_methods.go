package service

import (
	"errors"
	"net/url"
)

func (c *Client) Ping() error {
	return c.storage.Ping()
}

func (c *Client) CreateShortURL(originalURL, address string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}
	shortURL := c.genKey()
	err = c.storage.Put(shortURL, originalURL)
	if err != nil {
		return "", err
	}
	s, _ := url.JoinPath(address, shortURL)
	return s, nil
}

func (c *Client) GetShortURL(shortURL string) (string, error) {
	return c.storage.Get(shortURL)
}
