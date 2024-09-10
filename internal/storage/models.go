package storage

import "github.com/google/uuid"

type Row struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func newRow(shortURL, originalURL string) Row {
	return Row{
		UUID:        newUUID(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func newUUID() string {
	return uuid.NewString()
}
