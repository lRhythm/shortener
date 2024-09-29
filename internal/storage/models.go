package storage

import "github.com/google/uuid"

type Row struct {
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
	UserID        string `json:"user_id"`
	IsDeleted     bool   `json:"is_deleted"`
}

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

func newUUID() string {
	return uuid.NewString()
}
