package models

import "net/url"

type Row struct {
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url" db:"original_url"`
	CorrelationID string `json:"-"`
	IsDeleted     bool   `json:"-" db:"is_deleted"`
}

type Rows []Row

func (rs *Rows) ShortURLsWithAddress(address string) {
	for i, r := range *rs {
		s, _ := url.JoinPath(address, r.ShortURL)
		(*rs)[i].ShortURL = s
	}
}
