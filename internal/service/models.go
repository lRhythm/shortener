package service

import "github.com/lRhythm/shortener/internal/models"

// RepositoryInterface - интерфейс для имплементации хранилищем.
type RepositoryInterface interface {
	commonInterface
	CountURL() (cnt uint, err error)
	CountUser() (cnt uint, err error)
	Put(shortURL, originalURL, userID string) (err error)
	Batch(rows models.Rows, userID string) (err error)
	GetOriginalURL(shortURL string) (originalURL string, isDeleted bool, err error)
	GetShortURL(originalURL string) (shortURL string, err error)
	GetUserURLs(userID string) (rows models.Rows, err error)
	DeleteUserURLS(shortURLs []string, userID string) (err error)
}

// commonInterface - интерфейс вспомогательных методов хранилища.
type commonInterface interface {
	Ping() (err error)
	Close() (err error)
}

// Client - основной объект пакета для взаимодействия.
type Client struct {
	storage RepositoryInterface
}
