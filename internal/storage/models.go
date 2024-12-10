/*
Package storage - имплементации хранилища сервиса.
*/
package storage

import "github.com/google/uuid"

// Row - DTO сокращенного URL.
type Row struct {
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
	UserID        string `json:"user_id"`
	IsDeleted     bool   `json:"is_deleted"`
}

// newRow - конструктор Row.
func newRow(shortURL, originalURL, correlationID, userID string) Row {
	return Row{
		UUID:          newUUID(),
		CorrelationID: correlationID,
		ShortURL:      shortURL,
		OriginalURL:   originalURL,
		UserID:        userID,
		IsDeleted:     false,
	}
}

// newUUID - генератор UUID.
func newUUID() string {
	return uuid.NewString()
}
