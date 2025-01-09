package service

import (
	"errors"
	"net/url"

	"github.com/google/uuid"

	"github.com/lRhythm/shortener/internal/models"
)

// Ping - логика маршрута ping.
func (c *Client) Ping() error {
	return c.storage.Ping()
}

// CountURL - количество сокращённых URL в сервисе.
func (c *Client) CountURL() (cnt uint, err error) {
	return c.storage.CountURL()
}

// CountUser - количество пользователей в сервисе.
func (c *Client) CountUser() (cnt uint, err error) {
	return c.storage.CountUser()
}

// CreateShortURL - логика маршрута создания сокращенного URL.
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

// CreateBatch - логика маршрута пакетного создания сокращенного URL.
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

// GetOriginalURL - логика маршрута получения оригинального URL по сокращенному URL.
func (c *Client) GetOriginalURL(shortURL string) (string, bool, error) {
	return c.storage.GetOriginalURL(shortURL)
}

// GetUserURLs - логика маршрута получения сокращенных URL пользователя.
func (c *Client) GetUserURLs(address, userID string) (models.Rows, error) {
	rows, err := c.storage.GetUserURLs(userID)
	if err != nil {
		return nil, err
	}
	rows.ShortURLsWithAddress(address)
	return rows, nil
}

// DeleteUserURLs - логика маршрута удаления сокращенных URL пользователя.
func (c *Client) DeleteUserURLs(keys []string, userID string) {
	// Реализация Fan-In для соответствия требованиям.
	// Fan-In не требуется, т.к:
	// - в имплементации PostgreSQL удаление осуществляется с помощью 1 запроса;
	// - в имплементации InMemory сложность каждого вызова удаление 0n, где n - кол-во элементов слайса,
	inCh := genStrs(keys...)
	ch1 := pushStr(inCh)
	ch2 := pushStr(inCh)
	var values []string
	for n := range fanInStr(ch1, ch2) {
		values = append(values, n)
	}
	_ = c.storage.DeleteUserURLS(values, userID)
}

// GenerateUserID - логика генерации идентификатора пользователя в маршруте регистрации пользователя.
func (c *Client) GenerateUserID() string {
	return uuid.NewString()
}

// ValidateUserID - логика валидации идентификатора пользователя.
func (c *Client) ValidateUserID(userID string) error {
	_, err := uuid.Parse(userID)
	return err
}
