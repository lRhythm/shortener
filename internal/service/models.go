package service

import "github.com/lRhythm/shortener/internal/models"

type RepositoryInterface interface {
	Ping() (err error)
	Put(shortURL, originalURL string) (err error)
	Batch(rows models.Rows) (err error)
	GetOriginalURL(shortURL string) (originalURL string, err error)
	GetShortURL(originalURL string) (shortURL string, err error)
	Close() (err error)
}

type Client struct {
	storage RepositoryInterface
}
