/*
Package models - содержит объявление моделей, используемых в нескольких пакетах internal/.
*/
package models

import "net/url"

// Row - модель сокращенного URL.
type Row struct {
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url" db:"original_url"`
	CorrelationID string `json:"-"`
	IsDeleted     bool   `json:"-" db:"is_deleted"`
}

// Rows - коллекция моделей сокращенного URL.
type Rows []Row

// ShortURLsWithAddress - формирование полного адреса сокращенного URL для каждого элемента коллекции сокращенных URL.
func (rs *Rows) ShortURLsWithAddress(address string) {
	for i, r := range *rs {
		s, _ := url.JoinPath(address, r.ShortURL)
		(*rs)[i].ShortURL = s
	}
}
