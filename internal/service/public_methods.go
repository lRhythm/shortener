package service

import (
	"errors"
	"net/url"
)

func (c *Client) CreateShortURL(originalURL, address string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}
	keyURL := c.genKey()
	err = c.storage.Put(keyURL, originalURL)
	if err != nil {
		return "", err
	}
	s, _ := url.JoinPath(address, keyURL)
	return s, nil
}

func (c *Client) GetShortURL(keyURL string) (string, error) {
	return c.storage.Get(keyURL)
}
