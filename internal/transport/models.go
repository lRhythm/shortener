package transport

import "github.com/lRhythm/shortener/internal/models"

// ServiceInterface - интерфейс для имплементации сервисным слоем.
type ServiceInterface interface {
	CommonInterface
	URLInterface
	UserInterface
}

// CommonInterface - интерфейс вспомогательных методов.
type CommonInterface interface {
	Ping() (err error)
}

// URLInterface - интерфейс CRUD методов работы с URL.
type URLInterface interface {
	CountURL() (cnt uint, err error)
	CreateShortURL(originalURL, address, userID string) (shortURL string, err error)
	CreateBatch(rows models.Rows, address, userID string) (models.Rows, error)
	GetOriginalURL(key string) (originalURL string, isDeleted bool, err error)
	GetUserURLs(address, userID string) (rows models.Rows, err error)
	DeleteUserURLs(keys []string, userID string)
}

// UserInterface - интерфейс методов работы с пользователем.
type UserInterface interface {
	CountUser() (cnt uint, err error)
	GenerateUserID() string
	ValidateUserID(userID string) error
}

// CfgInterface - интерфейс методов работы с config.
type CfgInterface interface {
	Host() string
	Path() string
	TLSEnable() bool
	TLSPem() string
	TLSKey() string
	CookieKey() string
	Trusted() string
}
