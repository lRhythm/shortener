package service

import (
	"errors"
	"net/url"
)

func (c *Client) CreateShortURL(originalURL, host string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}
	keyURL := c.genKey()
	err = c.storage.Put(keyURL, originalURL)
	if err != nil {
		return "", err
	}
	s, _ := url.JoinPath(host, keyURL)
	return s, nil
}

func (c *Client) GetShortURL(keyURL string) (string, error) {
	return c.storage.Get(keyURL)
}
