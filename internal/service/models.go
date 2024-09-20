package service

import "github.com/lRhythm/shortener/internal/models"

type RepositoryInterface interface {
	commonInterface
	Put(shortURL, originalURL, userID string) (err error)
	Batch(rows models.Rows, userID string) (err error)
	GetOriginalURL(shortURL string) (originalURL string, err error)
	GetShortURL(originalURL string) (shortURL string, err error)
	GetUserURLs(userID string) (rows models.Rows, err error)
}

type commonInterface interface {
	Ping() (err error)
	Close() (err error)
}

type Client struct {
	storage RepositoryInterface
}
