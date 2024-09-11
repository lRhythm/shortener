package service

import (
	"errors"
	"github.com/lRhythm/shortener/internal/models"
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

func (c *Client) CreateBatch(rows models.Rows, address string) (models.Rows, error) {
	if len(rows) == 0 {
		return nil, errors.New("rows is empty")
	}
	for i, row := range rows {
		_, err := url.ParseRequestURI(row.OriginalURL)
		if err != nil {
			return nil, errors.New("invalid URL")
		}
		rows[i].ShortURL = c.genKey()
	}
	err := c.storage.Batch(rows)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		s, _ := url.JoinPath(address, row.ShortURL)
		rows[i].ShortURL = s
	}
	return rows, nil
}

func (c *Client) GetShortURL(shortURL string) (string, error) {
	return c.storage.Get(shortURL)
}
