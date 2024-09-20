package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lRhythm/shortener/internal/models"
	"net/url"
)

func (c *Client) Ping() error {
	return c.storage.Ping()
}

func (c *Client) CreateShortURL(originalURL, address, userID string) (string, error) {
	var err, e error
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}
	shortURL := c.genKey()
	err = c.storage.Put(shortURL, originalURL, userID)
	if err != nil {
		if !errors.Is(err, models.ErrConflict) {
			return "", err
		}
		shortURL, e = c.storage.GetShortURL(originalURL)
		if e != nil {
			return "", e
		}
		s, _ := url.JoinPath(address, shortURL)
		return s, err
	}
	s, _ := url.JoinPath(address, shortURL)
	return s, nil
}

func (c *Client) CreateBatch(rows models.Rows, address, userID string) (models.Rows, error) {
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
	err := c.storage.Batch(rows, userID)
	if err != nil {
		return nil, err
	}
	rows.ShortURLsWithAddress(address)
	return rows, nil
}

func (c *Client) GetOriginalURL(shortURL string) (string, error) {
	return c.storage.GetOriginalURL(shortURL)
}

func (c *Client) GetUserURLs(address, userID string) (models.Rows, error) {
	rows, err := c.storage.GetUserURLs(userID)
	if err != nil {
		return nil, err
	}
	rows.ShortURLsWithAddress(address)
	return rows, nil
}

func (c *Client) GenerateUserID() string {
	return uuid.NewString()
}

func (c *Client) ValidateUserID(userID string) error {
	_, err := uuid.Parse(userID)
	return err
}
