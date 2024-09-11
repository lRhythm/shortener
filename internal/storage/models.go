package storage

import "github.com/google/uuid"

type Row struct {
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

func newRow(shortURL, originalURL, correlationID string) Row {
	return Row{
		UUID:          newUUID(),
		CorrelationID: correlationID,
		ShortURL:      shortURL,
		OriginalURL:   originalURL,
	}
}

func newUUID() string {
	return uuid.NewString()
}
