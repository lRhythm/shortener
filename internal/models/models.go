package models

type Row struct {
	ShortURL      string
	OriginalURL   string
	CorrelationID string
}

type Rows []Row
